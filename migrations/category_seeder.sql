-- Category Seeder SQL
-- This file creates the nested category structure with proper lft/rgt values for nested set model

-- Insert root categories first
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
-- Concerts & Festivals (1-24)
('Concerts & Festivals', NULL, 'concerts_&_festivals', 'concerts-festivals', NULL, 1, 24, 0, NOW(), NOW()),

-- Party (25-42)
('Party', NULL, 'party', 'party', NULL, 25, 42, 1, NOW(), NOW()),

-- Culture (43-60)
('Culture', NULL, 'culture', 'culture', NULL, 43, 60, 2, NOW(), NOW()),

-- Shows (61-68)
('Shows', NULL, 'shows', 'shows', NULL, 61, 68, 3, NOW(), NOW()),

-- Sports (69-92)
('Sports', NULL, 'sports', 'sports', NULL, 69, 92, 4, NOW(), NOW()),

-- Freetime Activities (93-110)
('Freetime Activities', NULL, 'freetime_activities', 'freetime-activities', NULL, 93, 110, 5, NOW(), NOW()),

-- Business (111-126)
('Business', NULL, 'business', 'business', NULL, 111, 126, 6, NOW(), NOW()),

-- Ethnic (127-146)
('Ethnic', NULL, 'ethnic', 'ethnic', NULL, 127, 146, 7, NOW(), NOW()),

-- Other (147-150)
('Other', NULL, 'other', 'other', NULL, 147, 150, 8, NOW(), NOW());

-- Insert subcategories for Concerts & Festivals
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Rock & Pop', NULL, 'concerts_&_festivals', 'rock-pop', 1, 2, 3, 1, NOW(), NOW()),
('Hip-Hop & R&B', NULL, 'concerts_&_festivals', 'hip-hop-rb', 1, 4, 5, 1, NOW(), NOW()),
('Schlager & Volksmusic', NULL, 'concerts_&_festivals', 'schlager-volksmusic', 1, 6, 7, 1, NOW(), NOW()),
('Alternative & Indie Rock', NULL, 'concerts_&_festivals', 'alternative-indie-rock', 1, 8, 9, 1, NOW(), NOW()),
('Hard Rock & Metal', NULL, 'concerts_&_festivals', 'hard-rock-metal', 1, 10, 11, 1, NOW(), NOW()),
('Dance & Electro', NULL, 'concerts_&_festivals', 'dance-electro', 1, 12, 13, 1, NOW(), NOW()),
('Jazz & Blues', NULL, 'concerts_&_festivals', 'jazz-blues', 1, 14, 15, 1, NOW(), NOW()),
('Soul & Funk', NULL, 'concerts_&_festivals', 'soul-funk', 1, 16, 17, 1, NOW(), NOW()),
('Folk & Country', NULL, 'concerts_&_festivals', 'folk-country', 1, 18, 19, 1, NOW(), NOW()),
('World Music', NULL, 'concerts_&_festivals', 'world-music', 1, 20, 21, 1, NOW(), NOW()),
('K-Pop', NULL, 'concerts_&_festivals', 'k-pop', 1, 22, 23, 1, NOW(), NOW()),
('Other Concerts & Festivals', NULL, 'concerts_&_festivals', 'other-concerts-festivals', 1, 24, 25, 1, NOW(), NOW());

-- Insert subcategories for Party
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Bars', NULL, 'party', 'bars', 2, 26, 27, 1, NOW(), NOW()),
('Clubs', NULL, 'party', 'clubs', 2, 28, 29, 1, NOW(), NOW()),
('Rooftop Parties', NULL, 'party', 'rooftop-parties', 2, 30, 31, 1, NOW(), NOW()),
('Pub-Tours', NULL, 'party', 'pub-tours', 2, 32, 33, 1, NOW(), NOW()),
('Karakoke Nights', NULL, 'party', 'karakoke-nights', 2, 34, 35, 1, NOW(), NOW()),
('Exclusive VIP Events', NULL, 'party', 'exclusive-vip-events', 2, 36, 37, 1, NOW(), NOW()),
('Afterwork Party', NULL, 'party', 'afterwork-party', 2, 38, 39, 1, NOW(), NOW()),
('Other Parties', NULL, 'party', 'other-parties', 2, 40, 41, 1, NOW(), NOW());

-- Insert subcategories for Culture
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Theater', NULL, 'culture', 'theater', 3, 44, 45, 1, NOW(), NOW()),
('Musicals', NULL, 'culture', 'musicals', 3, 46, 47, 1, NOW(), NOW()),
('Classic Concerts', NULL, 'culture', 'classic-concerts', 3, 48, 49, 1, NOW(), NOW()),
('Opera & Operette', NULL, 'culture', 'opera-operette', 3, 50, 51, 1, NOW(), NOW()),
('Ballett & Dance', NULL, 'culture', 'ballett-dance', 3, 52, 53, 1, NOW(), NOW()),
('Art', NULL, 'culture', 'art', 3, 54, 55, 1, NOW(), NOW()),
('Fashion Shows', NULL, 'culture', 'fashion-shows', 3, 56, 57, 1, NOW(), NOW()),
('Culture Festivals', NULL, 'culture', 'culture-festivals', 3, 58, 59, 1, NOW(), NOW()),
('Cinema & Film', NULL, 'culture', 'cinema-film', 3, 60, 61, 1, NOW(), NOW());

-- Insert subcategories for Shows
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Dance Shows', NULL, 'shows', 'dance-shows', 4, 62, 63, 1, NOW(), NOW()),
('Comedy Shows', NULL, 'shows', 'comedy-shows', 4, 64, 65, 1, NOW(), NOW()),
('Kabarett', NULL, 'shows', 'kabarett', 4, 66, 67, 1, NOW(), NOW()),
('Iceshows', NULL, 'shows', 'iceshows', 4, 68, 69, 1, NOW(), NOW());

-- Insert subcategories for Sports
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Football', NULL, 'sports', 'football', 5, 70, 71, 1, NOW(), NOW()),
('American Football', NULL, 'sports', 'american-football', 5, 72, 73, 1, NOW(), NOW()),
('Basketball', NULL, 'sports', 'basketball', 5, 74, 75, 1, NOW(), NOW()),
('Boxing & Wrestling', NULL, 'sports', 'boxing-wrestling', 5, 76, 77, 1, NOW(), NOW()),
('Motor Sports', NULL, 'sports', 'motor-sports', 5, 78, 79, 1, NOW(), NOW()),
('Hockey', NULL, 'sports', 'hockey', 5, 80, 81, 1, NOW(), NOW()),
('Handball', NULL, 'sports', 'handball', 5, 82, 83, 1, NOW(), NOW()),
('Winter sports', NULL, 'sports', 'winter-sports', 5, 84, 85, 1, NOW(), NOW()),
('Tennis', NULL, 'sports', 'tennis', 5, 86, 87, 1, NOW(), NOW()),
('Ride Sports', NULL, 'sports', 'ride-sports', 5, 88, 89, 1, NOW(), NOW()),
('Fitness', NULL, 'sports', 'fitness', 5, 90, 91, 1, NOW(), NOW()),
('Other Sports', NULL, 'sports', 'other-sports', 5, 92, 93, 1, NOW(), NOW());

-- Insert subcategories for Freetime Activities
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Children/Kids', NULL, 'freetime_activities', 'children-kids', 6, 94, 95, 1, NOW(), NOW()),
('Activities (Indoor & Outdoor Activities)', NULL, 'freetime_activities', 'activities-indoor-outdoor-activities', 6, 96, 97, 1, NOW(), NOW()),
('Circus', NULL, 'freetime_activities', 'circus', 6, 98, 99, 1, NOW(), NOW()),
('Family', NULL, 'freetime_activities', 'family', 6, 100, 101, 1, NOW(), NOW()),
('Karneval', NULL, 'freetime_activities', 'karneval', 6, 102, 103, 1, NOW(), NOW()),
('Exhibitions (Ausstellungen)', NULL, 'freetime_activities', 'exhibitions-ausstellungen', 6, 104, 105, 1, NOW(), NOW()),
('Health & Wellness', NULL, 'freetime_activities', 'health-wellness', 6, 106, 107, 1, NOW(), NOW()),
('Charity Events', NULL, 'freetime_activities', 'charity-events', 6, 108, 109, 1, NOW(), NOW()),
('More Events', NULL, 'freetime_activities', 'more-events', 6, 110, 111, 1, NOW(), NOW());

-- Insert subcategories for Business
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Startups & Entrepreneurship', NULL, 'business', 'startups-entrepreneurship', 7, 112, 113, 1, NOW(), NOW()),
('Conferences & Summits', NULL, 'business', 'conferences-summits', 7, 114, 115, 1, NOW(), NOW()),
('Workshops and Training Sessions', NULL, 'business', 'workshops-training-sessions', 7, 116, 117, 1, NOW(), NOW()),
('Business Networking', NULL, 'business', 'business-networking', 7, 118, 119, 1, NOW(), NOW()),
('Career & Recruitment', NULL, 'business', 'career-recruitment', 7, 120, 121, 1, NOW(), NOW()),
('Trade Shows & Exhibitions', NULL, 'business', 'trade-shows-exhibitions', 7, 122, 123, 1, NOW(), NOW()),
('Corporate Events', NULL, 'business', 'corporate-events', 7, 124, 125, 1, NOW(), NOW()),
('Other Business Events', NULL, 'business', 'other-business-events', 7, 126, 127, 1, NOW(), NOW());

-- Insert subcategories for Ethnic
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Orient / Middle East', NULL, 'ethnic', 'orient-middle-east', 8, 128, 129, 1, NOW(), NOW()),
('Balkan', NULL, 'ethnic', 'balkan', 8, 130, 131, 1, NOW(), NOW()),
('Latin', NULL, 'ethnic', 'latin', 8, 132, 133, 1, NOW(), NOW()),
('European', NULL, 'ethnic', 'european', 8, 134, 135, 1, NOW(), NOW()),
('Arabic', NULL, 'ethnic', 'arabic', 8, 136, 137, 1, NOW(), NOW()),
('Asian', NULL, 'ethnic', 'asian', 8, 138, 139, 1, NOW(), NOW()),
('African', NULL, 'ethnic', 'african', 8, 140, 141, 1, NOW(), NOW()),
('USA', NULL, 'ethnic', 'usa', 8, 142, 143, 1, NOW(), NOW()),
('Indian', NULL, 'ethnic', 'indian', 8, 144, 145, 1, NOW(), NOW()),
('West Asian', NULL, 'ethnic', 'west-asian', 8, 146, 147, 1, NOW(), NOW());

-- Insert subcategories for Other
INSERT INTO categories (name, icon_id, type, slug, parent_id, lft, rgt, depth, created_at, updated_at) VALUES
('Miscellaneous', NULL, 'other', 'miscellaneous', 9, 148, 149, 1, NOW(), NOW()),
('Other Events', NULL, 'other', 'other-events', 9, 150, 151, 1, NOW(), NOW());