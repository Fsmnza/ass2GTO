ALTER TABLE module_info
    ADD CONSTRAINT check_updated_at CHECK (updated_at >= created_at),
    ADD CONSTRAINT check_module_duration CHECK (module_duration > 5 AND module_duration <= 15);