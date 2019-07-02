/* Ambrose-Daniel */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-f614f9c3670ad0475e819d76397abf0d',
  TRUE,
  'Ambrose-Daniel',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-7f48e2f2b6f7e4d1f9c864e48bc2b0f2',
  'c-f614f9c3670ad0475e819d76397abf0d'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-dc9537ca645ff34b4f289b6bd7aa08b7',
  'c-f614f9c3670ad0475e819d76397abf0d'
) ON CONFLICT DO NOTHING;

/* Ambrose-Isaac */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-d218888bdf510bbe1628d9983d75560f',
  TRUE,
  'Ambrose-Isaac',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-7f48e2f2b6f7e4d1f9c864e48bc2b0f2',
  'c-d218888bdf510bbe1628d9983d75560f'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-23e608245d0866ea937f15876adb5ef6',
  'c-d218888bdf510bbe1628d9983d75560f'
) ON CONFLICT DO NOTHING;

/* Ambrose-Sudharshan */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-fab2c2fb3befdbb2fe7abf444cbe3846',
  TRUE,
  'Ambrose-Sudharshan',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-7f48e2f2b6f7e4d1f9c864e48bc2b0f2',
  'c-fab2c2fb3befdbb2fe7abf444cbe3846'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-fb91825f564a3cc110f11836fedea6f4',
  'c-fab2c2fb3befdbb2fe7abf444cbe3846'
) ON CONFLICT DO NOTHING;

/* Daniel-Isaac */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-a1db4a9455dbc6c11ea2fa36f6bfa782',
  TRUE,
  'Daniel-Isaac',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-dc9537ca645ff34b4f289b6bd7aa08b7',
  'c-a1db4a9455dbc6c11ea2fa36f6bfa782'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-23e608245d0866ea937f15876adb5ef6',
  'c-a1db4a9455dbc6c11ea2fa36f6bfa782'
) ON CONFLICT DO NOTHING;

/* Daniel-Sudharshan */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-a3715860dcd95d1a105c12b7379e6d34',
  TRUE,
  'Daniel-Sudharshan',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-dc9537ca645ff34b4f289b6bd7aa08b7',
  'c-a3715860dcd95d1a105c12b7379e6d34'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-fb91825f564a3cc110f11836fedea6f4',
  'c-a3715860dcd95d1a105c12b7379e6d34'
) ON CONFLICT DO NOTHING;

/* Isaac-Sudharshan */
INSERT INTO "conversation" (
  "id", "dm", "title", "picture"
) VALUES (
  'c-6f2ba396fb53961ff8a6ba9c5d286a25',
  TRUE,
  'Isaac-Sudharshan',
  ''
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-23e608245d0866ea937f15876adb5ef6',
  'c-6f2ba396fb53961ff8a6ba9c5d286a25'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-fb91825f564a3cc110f11836fedea6f4',
  'c-6f2ba396fb53961ff8a6ba9c5d286a25'
) ON CONFLICT DO NOTHING;
