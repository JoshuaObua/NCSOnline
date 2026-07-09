-- Migration 054: Keep slideshow animation validation aligned with the CMS.

BEGIN;

ALTER TABLE cms_slides
    DROP CONSTRAINT IF EXISTS cms_slides_animation_type_check;

ALTER TABLE cms_slides
    ADD CONSTRAINT cms_slides_animation_type_check
    CHECK (animation_type IN (
        'fade-in',
        'slide-up',
        'slide-left',
        'slide-right',
        'zoom-in',
        'zoom-out',
        'flip-in',
        'blur-in',
        'bounce-in',
        'ken-burns'
    ));

COMMIT;
