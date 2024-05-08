-- 1. Простий запит на вибірку.
SELECT name FROM category;

-- 2. Запит на вибірку з використанням «between....and».
SELECT id, value FROM bid where value between 700 and 1000;

-- 3. Запит на вибірку з використанням «in».
SELECT id FROM category WHERE name IN ('Electronics', 'Automotive', 'Antiques');

-- 4. Запит на вибірку з використанням «like».
SELECT id, name, description FROM auction WHERE name LIKE '%Rare%';

-- 5. Запит на вибірку з двома умовами через «and».
SELECT id, name, is_active, created_at FROM auction WHERE is_active = true AND created_at > '05.06.2024';

-- 6. Запит на вибірку з двома умовами через «оr».
SELECT id, name, is_active, is_closed FROM auction WHERE is_active = false OR is_closed = true;

-- 7. Запит на вибірку з використанням «DISTINCT».
SELECT DISTINCT auction_lot_id FROM saved_auction_lots;

-- 8. Запит з функцією «min» або «max».
SELECT MAX(minimal_bid) FROM auction_lot;

-- 9. Запит з функцією «sum» або «avg».
SELECT AVG(bin_price) FROM auction_lot WHERE bin_price > 0;

-- 10. Запит з функцією «count».
SELECT COUNT(id) FROM auction_lot WHERE auction_id = 1;

-- 11. Запит на вибірку з використанням агрегатної функції і виведенням ще декількох полів.
SELECT a.id, a.name, AVG(b.value) FROM auction a
    INNER JOIN auction_lot l on a.id = l.auction_id
    INNER JOIN bid b on b.auction_lot_id = l.id
GROUP BY a.id ORDER BY a.id;

-- 12. Запит на вибірку з використанням агрегатної функції і умовою на вибірку поля.
SELECT id, name, AVG(bin_price) as avg_bin FROM auction_lot
    WHERE DATE(created_at) < '05.06.2024'
GROUP BY id;

-- 13. Запит на вибірку з використанням агрегатної функції і умовою на агрегатну функцію.
SELECT a.id, a.name, COUNT(l.id) as lot_count FROM auction a
    INNER JOIN auction_lot l on a.id = l.auction_id
GROUP BY a.id
HAVING COUNT(l.id) > 1;

-- 14. Запит на вибірку з використанням агрегатної функції, умовою на агрегатну функцію, умовою на вибірку поля з сортуванням даних.
SELECT a.id, a.name, AVG(l.minimal_bid) as avg_min_bid FROM auction a
    INNER JOIN auction_lot l on a.id = l.auction_id
WHERE l.reserve_price > 1000
GROUP BY a.id
HAVING AVG(minimal_bid) > 500;

-- 15. Запит з використанням INNER JOIN.
SELECT l.id, l.name, c.name as category_name FROM auction_lot l
INNER JOIN auction_lot_categories lc on lc.auction_lot_id = l.id
INNER JOIN category c on c.id = lc.category_id;

-- 16. Запит з використанням LEFT JOIN.
SELECT DISTINCT ON (b.id) b.id, b.value, b.auction_lot_id, l.is_active,
COALESCE((w.bid_id = b.id), false) AS did_win_lot FROM bid b
    LEFT JOIN auction_lot l on l.id = b.auction_lot_id
    LEFT JOIN auction_lot_winner w on w.bid_id = b.id
WHERE b.user_id = 5;

-- 17. Запит з використанням RIGHT JOIN.
SELECT b.id, b.value, w.won_at FROM bid b
RIGHT JOIN auction_lot_winner w on w.bid_id = b.id;

-- 18. Запит з використанням INNER JOIN і умовою.
SELECT l.id, l.name, c.name as category_name FROM auction_lot l
    INNER JOIN auction_lot_categories lc on lc.auction_lot_id = l.id
    INNER JOIN category c on c.id = lc.category_id
WHERE l.is_active = true AND l.is_closed = false;

-- 19. Запит з використанням INNER JOIN і умовою LIKE.
SELECT DISTINCT ON (a.id) a.id, a.name, a.is_active, a.owner_id, ct.name as category_name FROM auction a
    INNER JOIN auction_lot l on a.id = l.auction_id AND l.is_active = true
    INNER JOIN auction_lot_categories c on l.id = c.auction_lot_id
    INNER JOIN category ct on ct.id = c.category_id
WHERE a.name LIKE '%Rare%'
    AND c.category_id = 18
    AND a.is_active = true
    AND l.is_active = true;

-- 20. Запит з використанням INNER JOIN і використанням агрегатної функції.
SELECT l.id, l.name, AVG(b.value) as avg_bid FROM auction_lot l
    INNER JOIN bid b on b.auction_lot_id = l.id
GROUP BY l.id;

-- 21. Запит з використанням INNER JOIN і використанням агрегатної функції і умови HAVING.
SELECT a.id, a.name, AVG(l.minimal_bid) FROM auction a
    INNER JOIN auction_lot l on a.id = l.auction_id
GROUP BY a.id ORDER BY a.id
HAVING AVG(l.minimal_bid) > 100;

-- 22. Запит з використанням підзапита з використанням (=, <,>).
SELECT id, value, auction_lot_id FROM bid
WHERE value < (
    SELECT AVG(value) FROM bid
)
ORDER BY value desc LIMIT 5;

-- 23. Запит з використанням підзапита з використанням агрегатної функції.
SELECT id, value, auction_lot_id FROM bid
WHERE value > (
    SELECT AVG(value) FROM bid
);

-- 24. Запит з використанням підзапита з використанням оператора EXIST.
SELECT u.id, u.email
FROM users u
WHERE NOT EXISTS (
    SELECT 1
    FROM saved_auction_lots s
    WHERE s.user_id = u.id
);

-- 25. Запит з використанням підзапита з використанням АNY або SOME.
SELECT u.id, u.email FROM users u
WHERE u.id = ANY (
    SELECT user_id FROM bid
    WHERE DATE(created_at) = '05.06.2024'
    ORDER BY RANDOM() LIMIT 1
);

-- 26. Запит з використанням підзапита з використанням IN.
SELECT u.id, u.email FROM users u
WHERE u.id IN (
    SELECT user_id FROM bid ORDER BY value desc LIMIT 5
);

-- 27. Запит з використанням підзапита і зв’язку INNER JOIN.
SELECT l.id AS lot_id, l.name AS lot_name, b.value AS bid_value FROM auction_lot l
    INNER JOIN bid b ON b.auction_lot_id = l.id
GROUP BY l.id, b.value
HAVING b.value = (
    SELECT MAX(value) FROM bid WHERE bid.auction_lot_id = l.id
)
ORDER BY lot_id;
