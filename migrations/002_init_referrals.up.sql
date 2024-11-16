CREATE TABLE referral_codes (
    id SERIAL PRIMARY KEY,
    referrer_id INTEGER REFERENCES users(id),
    code VARCHAR(10) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_active_code_per_user UNIQUE (referrer_id, is_active) 
        WHERE is_active = true
);

CREATE TABLE referrals (
    id SERIAL PRIMARY KEY,
    referrer_id INTEGER REFERENCES users(id),
    referral_id INTEGER REFERENCES users(id),
    referral_code_id INTEGER REFERENCES referral_codes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_referral_codes_referrer ON referral_codes(referrer_id);
CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);