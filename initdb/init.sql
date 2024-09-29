CREATE TYPE address_type_enum AS ENUM ('Postal', 'Physical', 'Billing');

CREATE TABLE countries (
    country_code VARCHAR(3) PRIMARY KEY NOT NULL,
    country_name VARCHAR(255) NOT NULL
);


CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(3) REFERENCES countries(country_code),
    country VARCHAR(255),
    province VARCHAR(255),
    city VARCHAR(255),
    street3 VARCHAR(255),
    street2 VARCHAR(255),
    street VARCHAR(255),
    address_type address_type_enum,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE gender_enum AS ENUM ('Male', 'Female', 'Non-binary', 'Prefer not to say', 'Other');
CREATE TYPE status_enum AS ENUM ('Active', 'Inactive', 'Suspended');

------------------------------------------------------------- USERS
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    birthdate DATE NOT NULL,
    nationality VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    address_id BIGINT REFERENCES addresses(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    status status_enum DEFAULT 'Active',
    gender gender_enum,
    preferred_language VARCHAR(50),
    timezone VARCHAR(50),
    salt VARCHAR(255)
);

CREATE FUNCTION update_user_last_update()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_last_update();
------------------------------------------------------------- WORKFLOWS
CREATE TABLE workflows (
	id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    definition JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	author BIGINT REFERENCES users(id) NOT NULL
);

CREATE TABLE active_workflows (
	id SERIAL PRIMARY KEY NOT NULL,
	workflow BIGINT REFERENCES workflows(id) NOT NULL,
	current_state VARCHAR(255),
	date_submitted TIMESTAMPTZ,  -- Includes timezone information
	state_deadline TIMESTAMPTZ,  -- Includes timezone information
	workflow_instance JSONB NOT NULL
);

------------------------------------------------------------- TABLE DEFINITIONS
CREATE TYPE dispute_status AS ENUM (
    'Awaiting Respondant',
    'Active',
    'Review',
    'Settled',
    'Refused',
    'Withdrawn',
    'Transfer',
    'Appeal',
    'Other'
);

CREATE TABLE disputes (
	id SERIAL PRIMARY KEY,
	case_date DATE DEFAULT CURRENT_DATE,
	workflow BIGINT REFERENCES active_workflows(id) NOT NULL,
	status dispute_status DEFAULT 'Awaiting Respondant',
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	complainant BIGINT REFERENCES users(id) NOT NULL,
	respondant BIGINT REFERENCES users(id) NOT NULL,
    date_resolved DATE DEFAULT NULL
);

CREATE TABLE dispute_summaries (
	dispute BIGINT REFERENCES disputes(id) ON DELETE CASCADE,
	summary TEXT,
	PRIMARY KEY (dispute)
);

CREATE TABLE files (
	id SERIAL PRIMARY KEY,
	file_name VARCHAR(255) NOT NULL,
	uploaded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	file_path VARCHAR(255) NOT NULL
);

CREATE TABLE dispute_evidence (
	dispute BIGINT REFERENCES disputes(id),
	file_id BIGINT REFERENCES files(id) ON DELETE CASCADE NOT NULL,
	user_id BIGINT REFERENCES users(id) NOT NULL,
	PRIMARY KEY (dispute, file_id)
);

------------------------------------------------------------- TICKETING SYSTEM
CREATE TYPE ticket_status_enum AS ENUM (
    'Open',
    'Closed',
    'Solved',
    'On Hold'
);

CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  	-- Date the ticket was created
    created_by BIGINT REFERENCES users(id) NOT NULL, 	-- User that created the ticket
    dispute_id BIGINT REFERENCES disputes(id),       	-- Dispute the ticket is related to
    subject VARCHAR(255) NOT NULL,                   	-- Subject to describe the ticket
    status ticket_status_enum NOT NULL,              	-- Status of the ticket
    initial_message TEXT                             	-- Body of the initial message of the ticket
);

CREATE TABLE ticket_messages (
    id SERIAL PRIMARY KEY,
    ticket_id BIGINT REFERENCES tickets(id) NOT NULL,  	-- Reference to the ticket
    user_id BIGINT REFERENCES users(id) NOT NULL,      	-- User who made the comment
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    	-- Date the comment was submitted
    content TEXT NOT NULL                         		-- Body of the comment
);

------------------------------------------------------------- DISPUTE EXPERTS
CREATE TYPE expert_status AS ENUM ('Approved','Rejected','Review');

CREATE TABLE dispute_experts (
	dispute BIGINT REFERENCES disputes(id),
	"user" BIGINT REFERENCES users(id),
	PRIMARY KEY (dispute, "user")
);

CREATE TYPE exp_obj_status AS ENUM ('Review','Sustained','Overruled');

CREATE TABLE expert_objections (
    id SERIAL PRIMARY KEY,
    expert_id BIGINT REFERENCES users(id),                 -- Expert being objected to
    ticket_id BIGINT REFERENCES tickets(id) ON DELETE CASCADE,  -- Reference to the ticket
    status exp_obj_status DEFAULT 'Review'                 -- Status of the objection (Review, Sustained, Overruled)
);

-- View that automatically determines the status of the expert in the dispute
-- by examining the objections made to that expert
CREATE VIEW dispute_experts_view AS 
    SELECT dispute, "user" AS expert,
    (WITH statuses AS (
        SELECT eo.status FROM expert_objections eo
        JOIN tickets t ON eo.ticket_id = t.id
        WHERE t.dispute_id = dispute AND eo.expert_id = "user"
    ) SELECT CASE
        -- Expert is rejected if there exists a sustained objection
        WHEN 'Sustained' IN (SELECT * FROM statuses) THEN 'Rejected'::expert_status

        -- Expert is under review if there exist any objections that are under review
        WHEN 'Review' IN (SELECT * FROM statuses) THEN 'Review'::expert_status

        -- Approve the expert by default
        ELSE 'Approved'::expert_status
    END)
    AS status FROM dispute_experts;


CREATE FUNCTION check_valid_objection()
RETURNS trigger AS 
    $$
    DECLARE
        count integer;
    BEGIN
        -- Ensure the expert is assigned to the dispute referenced in the ticket
        SELECT COUNT(*) INTO count
        FROM dispute_experts de
        JOIN tickets t ON t.dispute_id = de.dispute
        WHERE de."user" = NEW.expert_id AND t.id = NEW.ticket_id;

        IF count = 0 THEN
            RAISE EXCEPTION 'Expert (ID = %) is not assigned to the dispute in ticket (ID = %)', NEW.expert_id, NEW.ticket_id;
        END IF;

        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER check_valid_objection
    BEFORE INSERT OR UPDATE ON expert_objections
    FOR EACH ROW
    EXECUTE FUNCTION check_valid_objection();


CREATE VIEW expert_objections_view AS
SELECT 
    eo.id AS objection_id,
	t.created_at AS objection_created_at,
    t.dispute_id,
    d.title AS dispute_title,
    eo.expert_id,
    expert.first_name || ' ' || expert.surname AS expert_full_name,
    t.created_by AS user_id,
    "user".first_name || ' ' || "user".surname AS user_full_name,
    t.initial_message AS reason,
    eo.status AS objection_status
FROM 
    expert_objections eo
JOIN 
    tickets t ON eo.ticket_id = t.id
JOIN 
    disputes d ON t.dispute_id = d.id
JOIN 
    users expert ON eo.expert_id = expert.id
JOIN 
    users "user" ON t.created_by = "user".id;


------------------------------------------------------------- EVENT LOG
CREATE TYPE event_types AS ENUM (
	'NOTIFICATION',
	'DISPUTE',
	'USER',
	'EXPERT',
	'WORKFLOW'
);

CREATE TABLE event_log (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	event_type event_types,
	event_data JSON
);


------------------------------------------------------------- TAGS
CREATE TABLE tags (
	id SERIAL PRIMARY KEY,
	tag_name VARCHAR(255) NOT NULL
);

CREATE TABLE dispute_tags (
	dispute_id BIGINT REFERENCES disputes(id) ON DELETE CASCADE,
	tag_id BIGINT REFERENCES tags(id),
	PRIMARY KEY (dispute_id, tag_id)
);

CREATE TABLE expert_tags (
	expert_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
	tag_id BIGINT REFERENCES tags(id),
	PRIMARY KEY (expert_id, tag_id)
);

CREATE TABLE workflow_tags (
	workflow_id BIGINT REFERENCES workflows(id) ON DELETE CASCADE,
	tag_id BIGINT REFERENCES tags(id),
	PRIMARY KEY (workflow_id, tag_id)
);
------------------------------------------------------------- DISPUTE DECISIONS
CREATE TABLE dispute_decisions (
    id SERIAL PRIMARY KEY,
    dispute_id BIGINT REFERENCES disputes(id),  					-- Reference to the dispute
    expert_id BIGINT REFERENCES users(id),      					-- Expert submitting the decision
    writeup_file_id BIGINT REFERENCES files(id),                  	-- Reference to the writeup file
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,             	-- Date the writeup was submitted
    UNIQUE (dispute_id)                                			  	-- One decision per dispute, regardless of who submitted it
);


------------------------------------------------------------- TABLE CONTENTS
INSERT INTO Countries (country_code, country_name) VALUES
('AF', 'Afghanistan'),
('AL', 'Albania'),
('DZ', 'Algeria'),
('AS', 'American Samoa'),
('AD', 'Andorra'),
('AO', 'Angola'),
('AI', 'Anguilla'),
('AQ', 'Antarctica'),
('AG', 'Antigua and Barbuda'),
('AR', 'Argentina'),
('AM', 'Armenia'),
('AW', 'Aruba'),
('AU', 'Australia'),
('AT', 'Austria'),
('AZ', 'Azerbaijan'),
('BS', 'Bahamas'),
('BH', 'Bahrain'),
('BD', 'Bangladesh'),
('BB', 'Barbados'),
('BY', 'Belarus'),
('BE', 'Belgium'),
('BZ', 'Belize'),
('BJ', 'Benin'),
('BM', 'Bermuda'),
('BT', 'Bhutan'),
('BO', 'Bolivia'),
('BA', 'Bosnia and Herzegovina'),
('BW', 'Botswana'),
('BV', 'Bouvet Island'),
('BR', 'Brazil'),
('IO', 'British Indian Ocean Territory'),
('BN', 'Brunei Darussalam'),
('BG', 'Bulgaria'),
('BF', 'Burkina Faso'),
('BI', 'Burundi'),
('KH', 'Cambodia'),
('CM', 'Cameroon'),
('CA', 'Canada'),
('CV', 'Cape Verde'),
('KY', 'Cayman Islands'),
('CF', 'Central African Republic'),
('TD', 'Chad'),
('CL', 'Chile'),
('CN', 'China'),
('CX', 'Christmas Island'),
('CC', 'Cocos (Keeling) Islands'),
('CO', 'Colombia'),
('KM', 'Comoros'),
('CG', 'Congo'),
('CD', 'Democratic Republic of the Congo'),
('CK', 'Cook Islands'),
('CR', 'Costa Rica'),
('CI', 'Côte d''Ivoire'),
('HR', 'Croatia'),
('CU', 'Cuba'),
('CY', 'Cyprus'),
('CZ', 'Czech Republic'),
('DK', 'Denmark'),
('DJ', 'Djibouti'),
('DM', 'Dominica'),
('DO', 'Dominican Republic'),
('EC', 'Ecuador'),
('EG', 'Egypt'),
('SV', 'El Salvador'),
('GQ', 'Equatorial Guinea'),
('ER', 'Eritrea'),
('EE', 'Estonia'),
('ET', 'Ethiopia'),
('FK', 'Falkland Islands (Malvinas)'),
('FO', 'Faroe Islands'),
('FJ', 'Fiji'),
('FI', 'Finland'),
('FR', 'France'),
('GF', 'French Guiana'),
('PF', 'French Polynesia'),
('TF', 'French Southern Territories'),
('GA', 'Gabon'),
('GM', 'Gambia'),
('GE', 'Georgia'),
('DE', 'Germany'),
('GH', 'Ghana'),
('GI', 'Gibraltar'),
('GR', 'Greece'),
('GL', 'Greenland'),
('GD', 'Grenada'),
('GP', 'Guadeloupe'),
('GU', 'Guam'),
('GT', 'Guatemala'),
('GG', 'Guernsey'),
('GN', 'Guinea'),
('GW', 'Guinea-Bissau'),
('GY', 'Guyana'),
('HT', 'Haiti'),
('HM', 'Heard Island and McDonald Islands'),
('VA', 'Holy See (Vatican City State)'),
('HN', 'Honduras'),
('HK', 'Hong Kong'),
('HU', 'Hungary'),
('IS', 'Iceland'),
('IN', 'India'),
('ID', 'Indonesia'),
('IR', 'Islamic Republic of Iran'),
('IQ', 'Iraq'),
('IE', 'Ireland'),
('IM', 'Isle of Man'),
('IL', 'Israel'),
('IT', 'Italy'),
('JM', 'Jamaica'),
('JP', 'Japan'),
('JE', 'Jersey'),
('JO', 'Jordan'),
('KZ', 'Kazakhstan'),
('KE', 'Kenya'),
('KI', 'Kiribati'),
('KP', 'Korea, Democratic People''s Republic of'),
('KR', 'Korea, Republic of'),
('KW', 'Kuwait'),
('KG', 'Kyrgyzstan'),
('LA', 'Lao People''s Democratic Republic'),
('LV', 'Latvia'),
('LB', 'Lebanon'),
('LS', 'Lesotho'),
('LR', 'Liberia'),
('LY', 'Libya'),
('LI', 'Liechtenstein'),
('LT', 'Lithuania'),
('LU', 'Luxembourg'),
('MO', 'Macao'),
('MK', 'Macedonia, the Former Yugoslav Republic of'),
('MG', 'Madagascar'),
('MW', 'Malawi'),
('MY', 'Malaysia'),
('MV', 'Maldives'),
('ML', 'Mali'),
('MT', 'Malta'),
('MH', 'Marshall Islands'),
('MQ', 'Martinique'),
('MR', 'Mauritania'),
('MU', 'Mauritius'),
('YT', 'Mayotte'),
('MX', 'Mexico'),
('FM', 'Micronesia, Federated States of'),
('MD', 'Moldova, Republic of'),
('MC', 'Monaco'),
('MN', 'Mongolia'),
('ME', 'Montenegro'),
('MS', 'Montserrat'),
('MA', 'Morocco'),
('MZ', 'Mozambique'),
('MM', 'Myanmar'),
('NA', 'Namibia'),
('NR', 'Nauru'),
('NP', 'Nepal'),
('NL', 'Netherlands'),
('NC', 'New Caledonia'),
('NZ', 'New Zealand'),
('NI', 'Nicaragua'),
('NE', 'Niger'),
('NG', 'Nigeria'),
('NU', 'Niue'),
('NF', 'Norfolk Island'),
('MP', 'Northern Mariana Islands'),
('NO', 'Norway'),
('OM', 'Oman'),
('PK', 'Pakistan'),
('PW', 'Palau'),
('PS', 'Palestinian Territory, Occupied'),
('PA', 'Panama'),
('PG', 'Papua New Guinea'),
('PY', 'Paraguay'),
('PE', 'Peru'),
('PH', 'Philippines'),
('PN', 'Pitcairn'),
('PL', 'Poland'),
('PT', 'Portugal'),
('PR', 'Puerto Rico'),
('QA', 'Qatar'),
('RE', 'Réunion'),
('RO', 'Romania'),
('RU', 'Russian Federation'),
('RW', 'Rwanda'),
('BL', 'Saint Barthélemy'),
('SH', 'Saint Helena, Ascension and Tristan da Cunha'),
('KN', 'Saint Kitts and Nevis'),
('LC', 'Saint Lucia'),
('MF', 'Saint Martin (French part)'),
('PM', 'Saint Pierre and Miquelon'),
('VC', 'Saint Vincent and the Grenadines'),
('WS', 'Samoa'),
('SM', 'San Marino'),
('ST', 'Sao Tome and Principe'),
('SA', 'Saudi Arabia'),
('SN', 'Senegal'),
('RS', 'Serbia'),
('SC', 'Seychelles'),
('SL', 'Sierra Leone'),
('SG', 'Singapore'),
('SX', 'Sint Maarten (Dutch part)'),
('SK', 'Slovakia'),
('SI', 'Slovenia'),
('SB', 'Solomon Islands'),
('SO', 'Somalia'),
('ZA', 'South Africa'),
('GS', 'South Georgia and the South Sandwich Islands'),
('SS', 'South Sudan'),
('ES', 'Spain'),
('LK', 'Sri Lanka'),
('SD', 'Sudan'),
('SR', 'Suriname'),
('SJ', 'Svalbard and Jan Mayen'),
('SZ', 'Swaziland'),
('SE', 'Sweden'),
('CH', 'Switzerland'),
('SY', 'Syrian Arab Republic'),
('TW', 'Taiwan, Province of China'),
('TJ', 'Tajikistan'),
('TZ', 'Tanzania, United Republic of'),
('TH', 'Thailand'),
('TL', 'Timor-Leste'),
('TG', 'Togo'),
('TK', 'Tokelau'),
('TO', 'Tonga'),
('TT', 'Trinidad and Tobago'),
('TN', 'Tunisia'),
('TR', 'Turkey'),
('TM', 'Turkmenistan'),
('TC', 'Turks and Caicos Islands'),
('TV', 'Tuvalu'),
('UG', 'Uganda'),
('UA', 'Ukraine'),
('AE', 'United Arab Emirates'),
('GB', 'United Kingdom'),
('US', 'United States'),
('UM', 'United States Minor Outlying Islands'),
('UY', 'Uruguay'),
('UZ', 'Uzbekistan'),
('VU', 'Vanuatu'),
('VE', 'Bolivarian Republic of Venezuela'),
('VN', 'Viet Nam'),
('VG', 'Virgin Islands, British'),
('VI', 'Virgin Islands, U.S.'),
('WF', 'Wallis and Futuna'),
('EH', 'Western Sahara'),
('YE', 'Yemen'),
('ZM', 'Zambia'),
('ZW', 'Zimbabwe');
