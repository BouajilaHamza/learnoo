import { test, expect } from '@playwright/test';

test('landing page works on mobile without JavaScript or external requests', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 390, height: 844 } });
  const page = await context.newPage();
  const external = [];
  page.on('request', request => { if (!request.url().startsWith('http://127.0.0.1:8080/')) external.push(request.url()); });
  await page.goto('http://127.0.0.1:8080/');
  await expect(page.getByRole('heading', { level: 1 })).toContainText('Le savoir se partage');
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link', { name: 'Aller au contenu' })).toBeFocused();
  await page.getByRole('link', { name: 'Découvrir le premier aperçu' }).click();
  await expect(page).toHaveURL(/#apprendre$/);
  await page.getByText('Est-ce que la plateforme est déjà disponible ?').click();
  await expect(page.locator('details').first()).toHaveAttribute('open', '');
  for (const width of [320, 390, 768, 1440]) {
    await page.setViewportSize({ width, height: 900 });
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    if (width === 390 || width === 1440) await page.screenshot({ path: `test-results/landing-${width}.png`, fullPage: true });
  }
  expect(external).toEqual([]);
  await page.getByRole('link', { name: "Ouvrir l'exemple technique non relu" }).click();
  await expect(page).toHaveURL(/\/lesson$/);
  await expect(page.getByRole('heading', { level: 1 })).toContainText('Nombres complexes');
  await context.close();
});
