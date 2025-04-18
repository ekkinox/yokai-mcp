-- +goose Up
CREATE TABLE books (
    id   INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    genre  VARCHAR(255) NOT NULL,
    synopsis VARCHAR(2048)
);

INSERT INTO books (title, genre, synopsis) VALUES
    ('The Silent Horizon', 'science-fiction', 'In a future where Earth’s atmosphere is slowly collapsing, a reclusive scientist must journey across a fractured world to deliver a formula that could save humanity. Battling distrust, rogue AI, and her own past, she discovers that the truth behind the disaster may be more terrifying than extinction.'),
    ('Beneath the Black Oak', 'horror', 'When a young widow returns to her ancestral home deep in the woods, she begins to unravel a legacy of madness and murder tied to an ancient, whispering tree. As she descends into a chilling spiral of hallucinations and family secrets, she must confront the darkness rooted both outside—and within.'),
    ('Chasing Light in Havana', 'romance', 'In 1950s Cuba, a headstrong photographer falls in love with a charming revolutionary. As political unrest builds, their passion grows, but so do the dangers. Torn between loyalty and love, she must choose between capturing the world as it is—or helping to change it forever.'),
    ('The Quantum Gambit', 'science-fiction', 'A disgraced quantum physicist is recruited by a covert agency to stop a rogue nation from activating a time-folding weapon. Racing against the clock through a maze of espionage and shifting timelines, he realizes that to stop the device, he may have to sacrifice the only reality he’s ever known.'),
    ('The Last Ember of Winter', 'fantasy', 'In a world locked in an eternal winter, an orphaned fire mage holds the key to restoring balance. Pursued by frostborn assassins and haunted by the memory of her burned village, she joins a band of outcasts in a perilous quest to reignite the world’s last ember—and bring back the sun.'),
    ('The Algorithm of Love', 'romance', 'A cynical app developer accidentally matches with his ex while testing a new dating algorithm. What starts as a prank turns into a journey through failed connections, awkward encounters, and rediscovered sparks—proving that love might just be the ultimate bug in the system.'),
    ('Tides of Mars', 'science-fiction', 'The Martian colonies are on the brink of war, and a disgraced pilot is thrust into the center of a rebellion when he rescues a fugitive scientist. With the fate of two planets hanging in the balance, alliances are tested, and legends are born in the crimson dust of the red planet.'),
    ('Letters to the Sky', 'romance', 'After a chance encounter during a delayed flight, a travel writer and a reserved aerospace engineer begin exchanging handwritten letters left in airport lounges around the world. As their bond deepens through stories, dreams, and confessions, they must decide if love can survive the leap from paper to real life.'),
    ('The Harvesting', 'horror', 'Every autumn, the townsfolk of Alder Hollow gather to celebrate the Harvest Festival—but this year, something is wrong. Crops bleed, scarecrows move, and people begin to vanish. When a skeptical reporter arrives to cover the quaint tradition, she uncovers an ancient pact between the town and a creature buried beneath the fields—one that demands its due.'),
    ('Crown of Feathers', 'fantasy', 'Born without magic in a kingdom ruled by spellcasters, a young stablehand discovers a hidden lineage tied to the last dragon of legend. As war looms and factions vie for power, she must risk everything to awaken the ancient beast—and claim a destiny written in fire and sky.');

-- +goose Down
DROP TABLE IF EXISTS books;