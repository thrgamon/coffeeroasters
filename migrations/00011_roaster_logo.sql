-- +goose Up
ALTER TABLE roasters ADD COLUMN logo_url VARCHAR(500);

UPDATE roasters SET logo_url = v.url FROM (VALUES
('1645-coffee', 'https://www.1645.com.au/cdn/shop/files/logo_200x200.gif?v=1614320773'),
('altura-coffee', 'https://alturacoffee.com.au/cdn/shop/files/Altura_Logo_Wordmark_Reverse_RGB.png?v=1715580991&width=300'),
('blackstar-coffee', 'https://blackstarcoffee.com.au/cdn/shop/files/BLACKSTAR-LOGO-2026-x01_600x.png?v=1772864929'),
('chico-loco', 'https://chico-loco-coffee-roasters-nt.myshopify.com/cdn/shop/files/logo_quality_f600ee48-6694-4d1b-a9f9-fe32547350db.png?v=1732401805&width=400'),
('dtown-coffee', 'https://dtowncoffeeroasters.com/cdn/shop/files/DTOWN-_Complete_Logo_-_Gold.png?v=1689739466&width=800'),
('fonzie-abbott', 'https://fonzieabbott.com/cdn/shop/files/FA-CircleBeerLogo-200x200px.png?v=1642654014'),
('little-marionette', 'https://thelittlemarionette.com/cdn/shop/files/TLM__logo_black_lg.png?v=1756179207'),
('market-lane', 'https://marketlane.com.au/cdn/shop/t/81/assets/logo.svg?v=101078261841729020861752936735'),
('mecca-coffee', 'https://www.mecca.coffee/cdn/shop/files/mecca_logo_-_Edited.webp?v=1756279867'),
('micrology', 'https://micrology.com.au/cdn/shop/files/micrology-checkout-logo.png?v=1637620963'),
('monastery-coffee', 'https://monastery.coffee/cdn/shop/files/MON001_Logo_FA-01_821c45c1-b1cd-4815-904b-ff03bc537a07_100x.png?v=1615319679'),
('ona-coffee', 'https://onacoffee.com.au/cdn/shop/files/logo-white_108x170.png?v=1625127507'),
('padre-coffee', 'https://www.padrecoffee.com.au/cdn/shop/files/Padre_crown_White.png?v=1768368848'),
('passport-coffee', 'https://passportcoffee.com.au/cdn/shop/files/passport-brisbane-logo_120x@2x.jpg?v=1664176689'),
('proud-mary', 'https://www.proudmarycoffee.com.au/cdn/shop/files/logo-new_cec94558-e7a4-450d-a83e-0ac32fb8694e_600x.png?v=1729739961'),
('seven-seeds', 'https://sevenseeds.com.au/cdn/shop/files/Seven-Seeds-Wordmark-Stacked.png?v=1723261709&width=160'),
('single-o', 'https://singleo.com.au/cdn/shop/files/New_Project_32x32.png?v=1690170011'),
('skittle-lane', 'https://skittlelane.com/cdn/shop/files/SL.svg?v=1731549873'),
('small-batch', 'https://www.smallbatch.com.au/wp-content/uploads/2013/05/logo-black-1024x1024.png'),
('soho-coffee', 'https://sohocoffeeroasters.com.au/cdn/shop/files/SOHO-LOGO-ROUND-1.jpg?v=1666661133&width=500'),
('twin-peaks', 'https://twinpeaks.net.au/cdn/shop/files/animated-logo002.gif?v=1769046138'),
('wolff-coffee', 'https://wolffcoffeeroasters.com.au/cdn/shop/files/wolff-logo-cin7_e06ff36a-db77-4eb1-940b-54c7d87a2103.png?v=1630599369')
) AS v(slug, url)
WHERE roasters.slug = v.slug;

-- +goose Down
ALTER TABLE roasters DROP COLUMN logo_url;
