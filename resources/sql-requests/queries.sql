-- 1. Простий запит на вибірку.
SELECT name FROM category;

-- 2. Запит на вибірку з використанням «between....and».
SELECT id, value FROM bids where value between 700 and 1000;

-- 3. Запит на вибірку з використанням «in».
SELECT * FROM app_user WHERE country IN ('usa', 'canada');

-- 4. Запит на вибірку з використанням «like».
SELECT * FROM auction WHERE name LIKE '%годинник%';

-- 5. Запит на вибірку з двома умовами через «and».
SELECT * FROM auction WHERE is_active = true AND created_at > '01.01.2024'

-- 6. Запит на вибірку з двома умовами через «оr».
SELECT * FROM auction WHERE is_active = false OR deleted_at IS NOT NULL;

-- 7. Запит на вибірку з використанням «DISTINCT».
SELECT DISTINCT(country) FROM app_user;

-- 8. Запит з функцією «min» або «max».
SELECT MAX(minimal_bid) FROM auction_lot;

-- 9. Запит з функцією «sum» або «avg».
SELECT AVG(bin_price) FROM auction_lot;

-- 10. Запит з функцією «count».
SELECT COUNT(id) FROM bids WHERE auction_lot_id = 1;

-- 11. Запит на вибірку з використанням агрегатної функції і виведенням ще декількох полів.
SELECT MAX(value), user_id FROM bids GROUP BY user_id;

-- 12. Запит на вибірку з використанням агрегатної функції і умовою на вибірку поля.
SELECT COUNT(id) FROM app_user WHERE country = 'usa';

-- 13. Запит на вибірку з використанням агрегатної функції і умовою на агрегатну функцію.


-- 14. Запит на вибірку з використанням агрегатної функції, умовою на агрегатну функцію, умовою на вибірку поля з сортуванням даних.


-- 15. Запит з використанням INNER JOIN.


-- 16. Запит з використанням LEFT JOIN.
SELECT a.id, a.name, a.created_at, u.email AS owner_email FROM auction a LEFT JOIN app_user u ON a.owner_id = u.id;

-- 17. Запит з використанням RIGHT JOIN.


-- 18. Запит з використанням INNER JOIN і умовою.









