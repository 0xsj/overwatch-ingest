"""PostgreSQL vector operations using pgvector extension."""

from enum import Enum
from typing import List, Optional, Tuple, Union

import numpy as np
from psycopg import Connection, Cursor, sql

from scout_common.database.postgres.client import Client
from scout_common.database.postgres.errors import map_error
from scout_common.errors import Error, invalid_field, validation_error


# Type alias for embeddings
Embedding = Union[List[float], np.ndarray]


class DistanceMetric(Enum):
    """
    Vector distance metrics supported by pgvector.
    
    - COSINE: Cosine distance (1 - cosine similarity), best for normalized vectors
    - L2: Euclidean distance, good for general purpose
    - INNER_PRODUCT: Negative inner product, useful for maximum inner product search
    """
    
    COSINE = "<->"           # Cosine distance operator
    L2 = "<->"               # L2 distance operator  
    INNER_PRODUCT = "<#>"    # Inner product operator


class IndexType(Enum):
    """
    Vector index types supported by pgvector.
    
    - IVFFLAT: Inverted file with flat compression (faster build, good recall)
    - HNSW: Hierarchical Navigable Small World (slower build, better recall)
    """
    
    IVFFLAT = "ivfflat"
    HNSW = "hnsw"


def ensure_extension(conn: Connection) -> None:
    """
    Ensure pgvector extension is installed and enabled.
    
    This should be called during database migrations or setup.
    
    Args:
        conn: Database connection
        
    Raises:
        Error: If extension cannot be created
        
    Example:
        with client.connection() as conn:
            ensure_extension(conn)
    """
    try:
        with conn.cursor() as cur:
            cur.execute("CREATE EXTENSION IF NOT EXISTS vector")
            conn.commit()
    except Exception as e:
        raise map_error(e, "create vector extension")


def normalize_embedding(embedding: Embedding) -> np.ndarray:
    """
    Normalize an embedding to unit length.
    
    This is important for cosine similarity searches as it allows
    using inner product instead of cosine distance for better performance.
    
    Args:
        embedding: Input embedding (list or numpy array)
        
    Returns:
        Normalized numpy array
        
    Example:
        embedding = [0.1, 0.2, 0.3]
        normalized = normalize_embedding(embedding)
    """
    arr = np.array(embedding, dtype=np.float32)
    norm = np.linalg.norm(arr)
    
    if norm == 0:
        return arr
    
    return arr / norm


def embedding_to_list(embedding: Embedding) -> List[float]:
    """
    Convert an embedding to a Python list for database storage.
    
    Args:
        embedding: Input embedding (list or numpy array)
        
    Returns:
        List of floats
        
    Example:
        embedding = np.array([0.1, 0.2, 0.3])
        db_format = embedding_to_list(embedding)
    """
    if isinstance(embedding, np.ndarray):
        return embedding.tolist()
    return list(embedding)


def validate_embedding_dimension(
    embedding: Embedding,
    expected_dim: int,
    field_name: str = "embedding",
) -> Optional[Error]:
    """
    Validate that an embedding has the expected dimensionality.
    
    Args:
        embedding: Embedding to validate
        expected_dim: Expected number of dimensions
        field_name: Field name for error message
        
    Returns:
        Error if validation fails, None otherwise
        
    Example:
        if err := validate_embedding_dimension(embedding, 1536, "text_embedding"):
            raise err
    """
    if isinstance(embedding, np.ndarray):
        actual_dim = embedding.shape[0]
    else:
        actual_dim = len(embedding)
    
    if actual_dim != expected_dim:
        return invalid_field(
            field_name,
            f"expected {expected_dim} dimensions, got {actual_dim}",
        )
    
    return None


def create_vector_column(
    conn: Connection,
    table_name: str,
    column_name: str,
    dimensions: int,
) -> None:
    """
    Add a vector column to an existing table.
    
    Args:
        conn: Database connection
        table_name: Name of the table
        column_name: Name of the vector column
        dimensions: Number of dimensions for the vector
        
    Raises:
        Error: If column creation fails
        
    Example:
        with client.connection() as conn:
            create_vector_column(conn, "incidents", "embedding", 1536)
    """
    try:
        with conn.cursor() as cur:
            query = sql.SQL(
                "ALTER TABLE {table} ADD COLUMN {column} vector({dims})"
            ).format(
                table=sql.Identifier(table_name),
                column=sql.Identifier(column_name),
                dims=sql.Literal(dimensions),
            )
            cur.execute(query)
            conn.commit()
    except Exception as e:
        raise map_error(e, f"create vector column {table_name}.{column_name}")


def create_vector_index(
    conn: Connection,
    table_name: str,
    column_name: str,
    index_type: IndexType = IndexType.HNSW,
    distance_metric: DistanceMetric = DistanceMetric.COSINE,
    lists: Optional[int] = None,
    m: Optional[int] = None,
    ef_construction: Optional[int] = None,
) -> None:
    """
    Create a vector index for efficient similarity search.
    
    Args:
        conn: Database connection
        table_name: Name of the table
        column_name: Name of the vector column
        index_type: Type of index (IVFFLAT or HNSW)
        distance_metric: Distance metric to use
        lists: Number of lists for IVFFLAT (default: rows/1000)
        m: Number of connections for HNSW (default: 16)
        ef_construction: Size of dynamic candidate list for HNSW (default: 64)
        
    Raises:
        Error: If index creation fails
        
    Example:
        # Create HNSW index with cosine distance
        with client.connection() as conn:
            create_vector_index(
                conn,
                "incidents",
                "embedding",
                index_type=IndexType.HNSW,
                distance_metric=DistanceMetric.COSINE,
            )
    """
    try:
        index_name = f"{table_name}_{column_name}_idx"
        
        if index_type == IndexType.IVFFLAT:
            # IVFFlat index
            # Default lists = rows / 1000 (will be calculated by postgres)
            lists_param = lists or 100
            
            with conn.cursor() as cur:
                query = sql.SQL(
                    "CREATE INDEX {index_name} ON {table} "
                    "USING ivfflat ({column} {metric}) "
                    "WITH (lists = {lists})"
                ).format(
                    index_name=sql.Identifier(index_name),
                    table=sql.Identifier(table_name),
                    column=sql.Identifier(column_name),
                    metric=sql.SQL(distance_metric.value),
                    lists=sql.Literal(lists_param),
                )
                cur.execute(query)
                conn.commit()
        
        elif index_type == IndexType.HNSW:
            # HNSW index
            m_param = m or 16
            ef_param = ef_construction or 64
            
            with conn.cursor() as cur:
                query = sql.SQL(
                    "CREATE INDEX {index_name} ON {table} "
                    "USING hnsw ({column} {metric}) "
                    "WITH (m = {m}, ef_construction = {ef})"
                ).format(
                    index_name=sql.Identifier(index_name),
                    table=sql.Identifier(table_name),
                    column=sql.Identifier(column_name),
                    metric=sql.SQL(distance_metric.value),
                    m=sql.Literal(m_param),
                    ef=sql.Literal(ef_param),
                )
                cur.execute(query)
                conn.commit()
        
    except Exception as e:
        raise map_error(e, f"create vector index on {table_name}.{column_name}")


def similarity_search(
    cur: Cursor,
    table_name: str,
    vector_column: str,
    query_embedding: Embedding,
    limit: int = 10,
    distance_metric: DistanceMetric = DistanceMetric.COSINE,
    where_clause: Optional[str] = None,
    where_params: Optional[tuple] = None,
    select_columns: Optional[List[str]] = None,
) -> List[Tuple]:
    """
    Perform similarity search using vector embeddings.
    
    Args:
        cur: Database cursor
        table_name: Name of the table to search
        vector_column: Name of the vector column
        query_embedding: Query embedding to search for
        limit: Maximum number of results to return
        distance_metric: Distance metric to use
        where_clause: Optional WHERE clause for filtering (without WHERE keyword)
        where_params: Parameters for WHERE clause
        select_columns: Columns to select (default: all columns + distance)
        
    Returns:
        List of tuples with search results
        
    Raises:
        Error: If search fails
        
    Example:
        # Simple similarity search
        with client.cursor() as cur:
            results = similarity_search(
                cur,
                "incidents",
                "embedding",
                query_embedding,
                limit=5,
            )
            
        # With filtering
        with client.cursor() as cur:
            results = similarity_search(
                cur,
                "incidents",
                "embedding",
                query_embedding,
                limit=5,
                where_clause="severity = %s AND status = %s",
                where_params=("HIGH", "OPEN"),
            )
    """
    try:
        # Convert embedding to list for database
        query_vector = embedding_to_list(query_embedding)
        
        # Build SELECT columns
        if select_columns:
            columns = ", ".join(select_columns)
        else:
            columns = "*"
        
        # Build distance expression
        distance_expr = f"{vector_column} {distance_metric.value} %s AS distance"
        
        # Build query
        query_parts = [
            f"SELECT {columns}, {distance_expr}",
            f"FROM {table_name}",
        ]
        
        params = [query_vector]
        
        if where_clause:
            query_parts.append(f"WHERE {where_clause}")
            if where_params:
                params.extend(where_params)
        
        query_parts.append(f"ORDER BY {vector_column} {distance_metric.value} %s")
        query_parts.append(f"LIMIT %s")
        
        params.append(query_vector)
        params.append(limit)
        
        query = " ".join(query_parts)
        
        cur.execute(query, tuple(params))
        return cur.fetchall()
    
    except Exception as e:
        raise map_error(e, f"similarity search on {table_name}")


def knn_search(
    cur: Cursor,
    table_name: str,
    vector_column: str,
    query_embedding: Embedding,
    k: int = 10,
    distance_metric: DistanceMetric = DistanceMetric.COSINE,
    distance_threshold: Optional[float] = None,
) -> List[Tuple]:
    """
    Perform k-nearest neighbors search.
    
    This is a specialized version of similarity_search optimized for KNN.
    
    Args:
        cur: Database cursor
        table_name: Name of the table to search
        vector_column: Name of the vector column
        query_embedding: Query embedding
        k: Number of nearest neighbors to return
        distance_metric: Distance metric to use
        distance_threshold: Optional maximum distance threshold
        
    Returns:
        List of tuples with (id, distance, ...)
        
    Raises:
        Error: If search fails
        
    Example:
        with client.cursor() as cur:
            neighbors = knn_search(
                cur,
                "incidents",
                "embedding",
                query_embedding,
                k=5,
                distance_threshold=0.3,  # Only return results within 0.3 distance
            )
    """
    where_clause = None
    where_params = None
    
    if distance_threshold is not None:
        where_clause = f"{vector_column} {distance_metric.value} %s < %s"
        where_params = (embedding_to_list(query_embedding), distance_threshold)
    
    return similarity_search(
        cur,
        table_name,
        vector_column,
        query_embedding,
        limit=k,
        distance_metric=distance_metric,
        where_clause=where_clause,
        where_params=where_params,
    )


def batch_insert_embeddings(
    cur: Cursor,
    table_name: str,
    vector_column: str,
    id_column: str,
    embeddings: List[Tuple[str, Embedding]],
) -> None:
    """
    Batch insert or update embeddings for multiple records.
    
    Args:
        cur: Database cursor
        table_name: Name of the table
        vector_column: Name of the vector column
        id_column: Name of the ID column
        embeddings: List of (id, embedding) tuples
        
    Raises:
        Error: If batch insert fails
        
    Example:
        embeddings = [
            ("incident-1", [0.1, 0.2, 0.3, ...]),
            ("incident-2", [0.4, 0.5, 0.6, ...]),
        ]
        
        with client.cursor() as cur:
            batch_insert_embeddings(
                cur,
                "incidents",
                "embedding",
                "id",
                embeddings,
            )
    """
    try:
        # Build UPSERT query
        query = sql.SQL(
            "INSERT INTO {table} ({id_col}, {vec_col}) "
            "VALUES (%s, %s) "
            "ON CONFLICT ({id_col}) "
            "DO UPDATE SET {vec_col} = EXCLUDED.{vec_col}"
        ).format(
            table=sql.Identifier(table_name),
            id_col=sql.Identifier(id_column),
            vec_col=sql.Identifier(vector_column),
        )
        
        # Convert embeddings to database format
        data = [
            (id_val, embedding_to_list(emb))
            for id_val, emb in embeddings
        ]
        
        cur.executemany(query, data)
    
    except Exception as e:
        raise map_error(e, f"batch insert embeddings into {table_name}")


def cosine_similarity(
    cur: Cursor,
    table_name: str,
    vector_column: str,
    id1: str,
    id2: str,
    id_column: str = "id",
) -> float:
    """
    Calculate cosine similarity between two stored embeddings.
    
    Args:
        cur: Database cursor
        table_name: Name of the table
        vector_column: Name of the vector column
        id1: ID of first record
        id2: ID of second record
        id_column: Name of ID column
        
    Returns:
        Cosine similarity (1 - cosine distance)
        
    Raises:
        Error: If calculation fails
        
    Example:
        with client.cursor() as cur:
            similarity = cosine_similarity(
                cur,
                "incidents",
                "embedding",
                "incident-1",
                "incident-2",
            )
            print(f"Similarity: {similarity}")
    """
    try:
        query = sql.SQL(
            "SELECT 1 - (a.{vec_col} <-> b.{vec_col}) "
            "FROM {table} a, {table} b "
            "WHERE a.{id_col} = %s AND b.{id_col} = %s"
        ).format(
            vec_col=sql.Identifier(vector_column),
            table=sql.Identifier(table_name),
            id_col=sql.Identifier(id_column),
        )
        
        cur.execute(query, (id1, id2))
        result = cur.fetchone()
        
        if result is None:
            raise validation_error(f"Could not find records with IDs {id1} or {id2}")
        
        return float(result[0])
    
    except Exception as e:
        raise map_error(e, f"calculate cosine similarity in {table_name}")