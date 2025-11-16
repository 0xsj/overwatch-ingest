"""Helper for creating vector-enabled migrations."""

from typing import Optional


def generate_vector_migration(
    table_name: str,
    column_name: str,
    dimensions: int,
    index_type: str = "hnsw",
    distance_metric: str = "cosine",
) -> str:
    """
    Generate SQL migration for adding vector support to a table.
    
    Args:
        table_name: Name of the table
        column_name: Name of the vector column
        dimensions: Number of dimensions
        index_type: Type of index (hnsw or ivfflat)
        distance_metric: Distance metric (cosine, l2, or inner_product)
        
    Returns:
        SQL migration string
        
    Example:
        sql = generate_vector_migration("incidents", "embedding", 1536)
        print(sql)
    """
    metric_map = {
        "cosine": "vector_cosine_ops",
        "l2": "vector_l2_ops",
        "inner_product": "vector_ip_ops",
    }
    
    ops_class = metric_map.get(distance_metric, "vector_cosine_ops")
    
    up_sql = f"""
-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Add vector column
ALTER TABLE {table_name} 
ADD COLUMN {column_name} vector({dimensions});

-- Create index
CREATE INDEX {table_name}_{column_name}_idx 
ON {table_name} 
USING {index_type} ({column_name} {ops_class});

-- Add comment
COMMENT ON COLUMN {table_name}.{column_name} 
IS 'Vector embedding ({dimensions} dimensions) for similarity search';
"""
    
    down_sql = f"""
-- Drop index
DROP INDEX IF EXISTS {table_name}_{column_name}_idx;

-- Drop column
ALTER TABLE {table_name} 
DROP COLUMN IF EXISTS {column_name};
"""
    
    return f"""-- Up Migration
{up_sql}

-- Down Migration
{down_sql}
"""