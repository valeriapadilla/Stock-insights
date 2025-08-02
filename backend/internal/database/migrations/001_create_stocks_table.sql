CREATE TABLE IF NOT EXISTS stocks (
    ticker TEXT NOT NULL,
    company TEXT NOT NULL,
    target_from TEXT,
    target_to TEXT,
    rating_from TEXT,
    rating_to TEXT,
    action TEXT,
    brokerage TEXT,
    time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (ticker, time)
);

CREATE INDEX IF NOT EXISTS idx_stocks_brokerage ON stocks(brokerage);
CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to);
CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC);
CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company);

CREATE INDEX IF NOT EXISTS idx_stocks_brokerage_time ON stocks(brokerage, time DESC);

CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);

COMMENT ON TABLE stocks IS 'Stocks data from external API'; 