-- Migration: 055_seed_mandate_history_pages
-- Seeds fully editable public Mandate and NCS History pages.

INSERT INTO cms_posts (
    id,
    title,
    slug,
    content,
    excerpt,
    category,
    category_tag,
    status,
    cover_image_url,
    author_id,
    published_at,
    meta_title,
    meta_description,
    focus_keywords,
    approved_at
)
VALUES
(
    'page_the_mandate',
    'Mandate',
    'the-mandate',
    ($mandate$
{
  "type": "ncs-page-builder",
  "version": 2,
  "blocks": [
    {
      "id": "mandate_intro",
      "type": "section",
      "props": {
        "background": "#ffffff",
        "padding": "0 0 20px",
        "margin": "0",
        "media": ""
      },
      "children": [
        {
          "id": "mandate_intro_heading",
          "type": "h2",
          "props": {
            "text": "Our Statutory Mandate",
            "align": "left"
          },
          "children": []
        },
        {
          "id": "mandate_intro_copy",
          "type": "paragraph",
          "props": {
            "text": "National Council of Sports (NCS) is a statutory organ whose establishment, status and powers are enshrined under the National Sports Act, 2023, to among other things, develop, promote and control sports activities in Uganda on behalf of Government, under the Ministry of Education and Sports. NCS is linked to the Supreme Council for Sports in Africa (SCSA) and other relevant sports organizations, and serves as a body corporate that coordinates all sports activities in the country in conjunction with National Sports Associations and Federations.",
            "lead": true
          },
          "children": []
        }
      ]
    },
    {
      "id": "mandate_body",
      "type": "sidebar_layout",
      "props": {
        "side": "right",
        "template": "minmax(0, 2fr) minmax(260px, 1fr)",
        "gap": "32px"
      },
      "children": [
        {
          "id": "mandate_main_column",
          "type": "column",
          "props": {
            "width": "2fr"
          },
          "children": [
            {
              "id": "mandate_functions_heading",
              "type": "h2",
              "props": {
                "text": "The Council Shall",
                "align": "left"
              },
              "children": []
            },
            {
              "id": "mandate_functions_list",
              "type": "list",
              "props": {
                "items": "Recognize a sports discipline as a national sports discipline.\nRegister national sports organisations.\nPromote and regulate the activities of national sports associations and national sports federations and, where necessary, award medals, certificates of recognition, trophies and other incentives.\nIn collaboration with national sports associations, national sports federations, local governments, educational institutions, communities and the private sector, make provision for sports facilities, equipment and training; promote sportsmanship by searching for, identifying and developing sporting talent and ensuring discipline among sportspersons; and create public awareness through sporting events on matters of national interest and the benefits of sports to health.\nOrganise sports clinics and provide advisory and counselling services to athletes.\nDevelop, manage, operate and maintain the public sports facilities vested in the Council under the Act.\nEstablish, operate and maintain sports museums.\nApprove the expenditure by national sports associations and national sports federations of funds and grants received from Government.\nFacilitate cooperation between and among national sports associations and national sports federations.\nIn collaboration with the Ministry, facilitate the participation of Ugandan athletes and national teams in international sports competitions.\nApprove the hosting of international sports competitions and sports festivals by national sports associations and national sports federations.\nPerform any other function as may be required under the Act.",
                "ordered": true
              },
              "children": []
            },
            {
              "id": "mandate_additional_card",
              "type": "card",
              "props": {
                "title": "Additional Statutory Responsibilities",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "mandate_additional_list",
                  "type": "list",
                  "props": {
                    "items": "Maintain registers of the national sports organisations registered under the NCS Act in a format prescribed by regulations made under the Act.\nCooperate with the Ministry and other Government ministries, departments and agencies in the implementation of the NCS Act.",
                    "ordered": false
                  },
                  "children": []
                }
              ]
            }
          ]
        },
        {
          "id": "mandate_sidebar",
          "type": "column",
          "props": {
            "width": "1fr"
          },
          "children": [
            {
              "id": "mandate_vision_card",
              "type": "card",
              "props": {
                "title": "Vision",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "mandate_vision_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "A centre of excellence for promotion and development of Sports.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "mandate_mission_card",
              "type": "card",
              "props": {
                "title": "Mission",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "mandate_mission_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "Maximizing opportunities for all Ugandans to participate and excel in Sports.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "mandate_legal_card",
              "type": "card",
              "props": {
                "title": "Legal Foundation",
                "footer": "National Sports Act, 2023",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "mandate_legal_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "NCS operates under the Ministry of Education and Sports as Uganda's statutory body for the development, promotion and regulation of sports.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "id": "mandate_values_section",
      "type": "section",
      "props": {
        "background": "#f8fafc",
        "padding": "40px",
        "margin": "28px 0 0",
        "media": ""
      },
      "children": [
        {
          "id": "mandate_values_heading",
          "type": "h2",
          "props": {
            "text": "Guiding Principles",
            "align": "left"
          },
          "children": []
        },
        {
          "id": "mandate_values_intro",
          "type": "paragraph",
          "props": {
            "text": "National Council of Sports is guided by five core values.",
            "lead": true
          },
          "children": []
        },
        {
          "id": "mandate_values_grid",
          "type": "grid",
          "props": {
            "template": "repeat(2, minmax(0, 1fr))",
            "gap": "18px"
          },
          "children": [
            {
              "id": "mandate_values_column_one",
              "type": "column",
              "props": {
                "width": "1fr"
              },
              "children": [
                {
                  "id": "mandate_honesty_card",
                  "type": "card",
                  "props": {
                    "title": "1. Honesty",
                    "footer": "",
                    "padding": "20px"
                  },
                  "children": [
                    {
                      "id": "mandate_honesty_copy",
                      "type": "paragraph",
                      "props": {
                        "text": "Moral character that implies positive attributes such as truthfulness, integrity, being straightforward and not cheating. It also involves being sincere, loyal, trustworthy and fair.",
                        "lead": false
                      },
                      "children": []
                    }
                  ]
                },
                {
                  "id": "mandate_excellence_card",
                  "type": "card",
                  "props": {
                    "title": "2. Pursuit of Personal Excellence",
                    "footer": "",
                    "padding": "20px"
                  },
                  "children": [
                    {
                      "id": "mandate_excellence_copy",
                      "type": "paragraph",
                      "props": {
                        "text": "Positive change to fulfil your dreams through hard work, building self-confidence, breaking long-term goals into manageable tasks, learning from the best and maintaining a strong desire to succeed.",
                        "lead": false
                      },
                      "children": []
                    }
                  ]
                },
                {
                  "id": "mandate_sport_card",
                  "type": "card",
                  "props": {
                    "title": "3. Love of Sport",
                    "footer": "",
                    "padding": "20px"
                  },
                  "children": [
                    {
                      "id": "mandate_sport_copy",
                      "type": "paragraph",
                      "props": {
                        "text": "A feeling of attachment to sports. Sport supports mental health and physical fitness and facilitates the attainment of personal, community and national objectives.",
                        "lead": false
                      },
                      "children": []
                    }
                  ]
                }
              ]
            },
            {
              "id": "mandate_values_column_two",
              "type": "column",
              "props": {
                "width": "1fr"
              },
              "children": [
                {
                  "id": "mandate_teamwork_card",
                  "type": "card",
                  "props": {
                    "title": "4. Team Work",
                    "footer": "",
                    "padding": "20px"
                  },
                  "children": [
                    {
                      "id": "mandate_teamwork_copy",
                      "type": "paragraph",
                      "props": {
                        "text": "Supporting each other to drive the Council to high performance.",
                        "lead": false
                      },
                      "children": []
                    }
                  ]
                },
                {
                  "id": "mandate_inclusiveness_card",
                  "type": "card",
                  "props": {
                    "title": "5. Inclusiveness",
                    "footer": "",
                    "padding": "20px"
                  },
                  "children": [
                    {
                      "id": "mandate_inclusiveness_copy",
                      "type": "paragraph",
                      "props": {
                        "text": "Providing equal access to opportunities and resources for people who might otherwise be excluded or marginalised, including people with physical or mental disabilities and members of minority groups.",
                        "lead": false
                      },
                      "children": []
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
$mandate$::jsonb)::text,
    'The statutory mandate, vision, mission and guiding principles of the National Council of Sports.',
    'page',
    'institutional',
    'published',
    '',
    NULL,
    NOW(),
    'Mandate | National Council of Sports Uganda',
    'Read the statutory mandate, vision, mission and guiding principles of the National Council of Sports Uganda.',
    'NCS mandate, National Sports Act 2023, sports regulation Uganda',
    NOW()
),
(
    'page_ncs_history',
    'NCS History',
    'ncs-history',
    ($history$
{
  "type": "ncs-page-builder",
  "version": 2,
  "blocks": [
    {
      "id": "history_intro",
      "type": "section",
      "props": {
        "background": "#ffffff",
        "padding": "0 0 20px",
        "margin": "0",
        "media": ""
      },
      "children": [
        {
          "id": "history_intro_heading",
          "type": "h2",
          "props": {
            "text": "More Than Six Decades of Service to Ugandan Sport",
            "align": "left"
          },
          "children": []
        },
        {
          "id": "history_intro_copy",
          "type": "paragraph",
          "props": {
            "text": "The National Council of Sports has served as Uganda's statutory national sports body since 1964. Across the decades, the Council has coordinated sports development, supported national sports organisations, guided participation in international competition and helped create opportunities for Ugandans to participate and excel in sport.",
            "lead": true
          },
          "children": []
        }
      ]
    },
    {
      "id": "history_body",
      "type": "sidebar_layout",
      "props": {
        "side": "right",
        "template": "minmax(0, 2fr) minmax(260px, 1fr)",
        "gap": "32px"
      },
      "children": [
        {
          "id": "history_main_column",
          "type": "column",
          "props": {
            "width": "2fr"
          },
          "children": [
            {
              "id": "history_timeline_heading",
              "type": "h2",
              "props": {
                "text": "Our Journey",
                "align": "left"
              },
              "children": []
            },
            {
              "id": "history_1964_card",
              "type": "card",
              "props": {
                "title": "1964 - The Council Is Established",
                "footer": "National Council of Sports Act, Chapter 48",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_1964_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "The National Council of Sports was established under the National Council of Sports Act. The Act was assented to on 22 June 1964 and commenced on 25 June 1964, creating a statutory national body to develop, promote and coordinate sport in Uganda.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "history_service_card",
              "type": "card",
              "props": {
                "title": "A National Coordinating Role",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_service_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "NCS developed its role as the apex body coordinating sports activities across the country. Working with national sports associations and federations, local governments, educational institutions, communities and partners, the Council supported facilities, training, talent development and national representation.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "history_2023_card",
              "type": "card",
              "props": {
                "title": "2023 - A Modern National Sports Framework",
                "footer": "National Sports Act, 2023",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_2023_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "The National Sports Act, 2023 modernised Uganda's sports governance framework and reaffirmed the Council's statutory responsibilities, including recognising sports disciplines, registering national sports organisations, regulating their activities and facilitating Ugandan participation in international competition.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "history_today_card",
              "type": "card",
              "props": {
                "title": "NCS Today",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_today_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "Today, NCS continues to operate under the Ministry of Education and Sports, connecting Government, sports organisations, athletes, coaches, administrators and communities around a shared goal: a strong, inclusive and high-performing sports sector for Uganda.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            }
          ]
        },
        {
          "id": "history_sidebar",
          "type": "column",
          "props": {
            "width": "1fr"
          },
          "children": [
            {
              "id": "history_dates_card",
              "type": "card",
              "props": {
                "title": "At a Glance",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_dates_list",
                  "type": "list",
                  "props": {
                    "items": "22 June 1964: founding Act assented to.\n25 June 1964: founding Act commenced.\n2023: modern national sports legislation enacted.\nPresent: NCS continues to lead national sports development.",
                    "ordered": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "history_vision_card",
              "type": "card",
              "props": {
                "title": "Our Continuing Vision",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_vision_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "A centre of excellence for promotion and development of Sports.",
                    "lead": false
                  },
                  "children": []
                }
              ]
            },
            {
              "id": "history_mandate_card",
              "type": "card",
              "props": {
                "title": "Our Work Today",
                "footer": "",
                "padding": "20px"
              },
              "children": [
                {
                  "id": "history_mandate_copy",
                  "type": "paragraph",
                  "props": {
                    "text": "Explore the Council's current responsibilities under Uganda's national sports framework.",
                    "lead": false
                  },
                  "children": []
                },
                {
                  "id": "history_mandate_link",
                  "type": "button_primary",
                  "props": {
                    "label": "Read the Mandate",
                    "href": "/pages/the-mandate",
                    "variant": "primary",
                    "icon": "icofont-rounded-right",
                    "iconPosition": "right"
                  },
                  "children": []
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "id": "history_legacy_section",
      "type": "section",
      "props": {
        "background": "#f8fafc",
        "padding": "40px",
        "margin": "28px 0 0",
        "media": ""
      },
      "children": [
        {
          "id": "history_legacy_heading",
          "type": "h2",
          "props": {
            "text": "A Legacy Built Through Partnership",
            "align": "left"
          },
          "children": []
        },
        {
          "id": "history_legacy_copy",
          "type": "paragraph",
          "props": {
            "text": "The history of NCS is also the history of the athletes, coaches, volunteers, administrators, associations, federations, institutions and communities that have advanced Ugandan sport. The Council's continuing role is to bring these efforts together, protect standards and expand opportunities for future generations.",
            "lead": true
          },
          "children": []
        }
      ]
    }
  ]
}
$history$::jsonb)::text,
    'The history of the National Council of Sports, from its establishment in 1964 to its work today.',
    'page',
    'institutional',
    'published',
    '',
    NULL,
    NOW(),
    'NCS History | National Council of Sports Uganda',
    'Explore the history of the National Council of Sports Uganda from its establishment in 1964 to the present day.',
    'NCS history, Uganda sports history, National Council of Sports 1964',
    NOW()
)
ON CONFLICT (slug) DO UPDATE
SET title = EXCLUDED.title,
    content = EXCLUDED.content,
    excerpt = EXCLUDED.excerpt,
    category = EXCLUDED.category,
    category_tag = EXCLUDED.category_tag,
    status = EXCLUDED.status,
    published_at = COALESCE(cms_posts.published_at, EXCLUDED.published_at),
    meta_title = EXCLUDED.meta_title,
    meta_description = EXCLUDED.meta_description,
    focus_keywords = EXCLUDED.focus_keywords,
    approved_at = COALESCE(cms_posts.approved_at, EXCLUDED.approved_at),
    updated_at = NOW();

DO $$
DECLARE
    seeded_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO seeded_count
    FROM cms_posts
    WHERE slug IN ('the-mandate', 'ncs-history')
      AND category = 'page'
      AND status = 'published'
      AND content::jsonb ->> 'type' = 'ncs-page-builder';

    IF seeded_count <> 2 THEN
        RAISE EXCEPTION 'Expected two published page-builder pages, found %', seeded_count;
    END IF;
END $$;
