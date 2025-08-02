CREATE TABLE IF NOT EXISTS recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT NOT NULL,
    score NUMERIC NOT NULL,
    explanation TEXT,
    run_at TIMESTAMPTZ NOT NULL,
    rank INTEGER
);

CREATE INDEX IF NOT EXISTS idx_recommendations_run_at ON recommendations(run_at DESC);
CREATE INDEX IF NOT EXISTS idx_recommendations_score ON recommendations(score DESC);
CREATE INDEX IF NOT EXISTS idx_recommendations_rank ON recommendations(rank);
CREATE INDEX IF NOT EXISTS idx_recommendations_ticker ON recommendations(ticker);

-- Composite index for recommendation queries
CREATE INDEX IF NOT EXISTS idx_recommendations_run_at_score ON recommendations(run_at DESC, score DESC);

-- Add comment to table
COMMENT ON TABLE recommendations IS 'Stock recommendations calculated by algorithm'; 