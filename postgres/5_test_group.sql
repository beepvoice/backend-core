INSERT INTO "conversation" (
  "id", "title"
) VALUES (
  'c-d73b6afa2fe3685faad28eba36d8cd0a',
  'AWESOME'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-7f48e2f2b6f7e4d1f9c864e48bc2b0f2',
  'c-d73b6afa2fe3685faad28eba36d8cd0a'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-dc9537ca645ff34b4f289b6bd7aa08b7',
  'c-d73b6afa2fe3685faad28eba36d8cd0a'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-23e608245d0866ea937f15876adb5ef6',
  'c-d73b6afa2fe3685faad28eba36d8cd0a'
) ON CONFLICT DO NOTHING;

INSERT INTO "member" (
  "user", "conversation"
) VALUES (
  'u-fb91825f564a3cc110f11836fedea6f4',
  'c-d73b6afa2fe3685faad28eba36d8cd0a'
) ON CONFLICT DO NOTHING;
