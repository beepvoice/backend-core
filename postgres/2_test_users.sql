INSERT INTO "user" (
  id, first_name, last_name, phone_number
) VALUES (
  'u-7f48e2f2b6f7e4d1f9c864e48bc2b0f2',
  'Ambrose',
  'Chua',
  '+65 9766 3827'
) ON CONFLICT DO NOTHING;

INSERT INTO "user" (
  id, first_name, last_name, phone_number
) VALUES (
  'u-dc9537ca645ff34b4f289b6bd7aa08b7',
  'Daniel',
  'Lim',
  '+65 8737 7117'
) ON CONFLICT DO NOTHING;

INSERT INTO "user" (
  id, first_name, last_name, phone_number
) VALUES (
  'u-23e608245d0866ea937f15876adb5ef6',
  'Isaac',
  'Tay',
  '+65 8181 6346'
) ON CONFLICT DO NOTHING;

INSERT INTO "user" (
  id, first_name, last_name, phone_number
) VALUES (
  'u-fb91825f564a3cc110f11836fedea6f4',
  'Sudharshan',
  '',
  '+65 8143 8417'
) ON CONFLICT DO NOTHING;
