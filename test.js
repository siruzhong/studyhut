'use strict';
const puppeteer = require('puppeteer');
async function test() {
    const browser = await puppeteer.launch({args: ['--no-sandbox', '--disable-setuid-sandbox'], headless: true, ignoreHTTPSErrors: true});
    const page = await browser.newPage();
    await page.goto('http://localhost:8181', {"waitUntil" :  ['networkidle2', 'domcontentloaded'], "timeout": 10000});
    let content=await page.content();
    console.log(content);
    await browser.close();
}
test();