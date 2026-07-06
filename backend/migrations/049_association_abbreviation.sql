-- Migration: 049_association_abbreviation
-- Adds a short official abbreviation (e.g. "FUFA") to associations, shown
-- as a badge next to the name on the public Associations directory —
-- distinct from an initialism derived from the name, since real federation
-- abbreviations don't always match first-letter acronyms.

ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS abbreviation TEXT NOT NULL DEFAULT '';
