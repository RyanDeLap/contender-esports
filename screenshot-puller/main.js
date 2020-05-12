//const puppeteer = require('puppeteer-core');
const puppeteer = require('puppeteer');
const atob = require('atob');
const btoa = require('btoa');
const fs = require('fs');

//var url = 'https://slider.ggleap.com/?center=476fc5ba-6114-4445-94f8-b3734e7f770d&screen=main';
var wantedCategory = "Fortnite";
var wantedTimeframe = "current_week";

// ty google
function evaluate(page, func) {
    var args = [].slice.call(arguments, 2);
    var fn = "(function() { return (" + func.toString() + ").apply(this, " + JSON.stringify(args) + ");})();";
    return page.evaluate(fn);
}

async function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function run() {
    let arguments = process.argv.slice(2);
    if(arguments.length < 2) {
        console.log("Error, Command needed: node main.js game filename url")
        process.exit(1)
    }
    let game = arguments[0]
    let filename = arguments[1]
    let url = arguments[2]

    console.log("Starting image puller. Getting " + game + " for url " + url)

    //const browser = await puppeteer.connect({ browserWSEndpoint: 'ws://172.s17.0.2:3001' });
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    const client = await page.target().createCDPSession();
    await client.send('Network.enable');
    await client.send('Network.setRequestInterception', {
        patterns: [
            {
                urlPattern: '*',
                resourceType: 'Script',
                interceptionStage: 'HeadersReceived'
            }
        ]
    });
    client.on('Network.requestIntercepted', async ({ interceptionId, request, responseHeaders, resourceType }) => {
        console.log(`Intercepted ${request.url} {interception id: ${interceptionId}}`);

        const response = await client.send('Network.getResponseBodyForInterception', {
            interceptionId
        });
        const originalBody = response.base64Encoded ? atob(response.body) : response.body;
        const contentTypeHeader = Object.keys(responseHeaders).find(k => k.toLowerCase() === 'content-type');
        let newBody, contentType = responseHeaders[contentTypeHeader];
        newBody = originalBody

        if (request.url.includes("main.") && request.url.endsWith(".js")) {
            newBody = newBody.replace(/,([^\d]+?)\.prototype\.showNextSlide=/, function(useless, group1){
                return ",window.apiExposer=" + group1 + ","+group1+".prototype.showNextSlide=";
            });
        }

        const newHeaders = [
            'Date: ' + (new Date()).toUTCString(),
            'Connection: closed',
            'Content-Length: ' + newBody.length,
            'Content-Type: ' + contentType
        ];

        console.log(`Continuing interception ${interceptionId}`)
        client.send('Network.continueInterceptedRequest', {
            interceptionId,
            rawResponse: btoa('HTTP/1.1 200 OK' + '\r\n' + newHeaders.join('\r\n') + '\r\n\r\n' + newBody)
        });
    });

    page.on('console', consoleObj => console.log(consoleObj.text()));

    await page.goto(url);
    await page.waitForSelector("body")
    await page.setViewport({
        width: 1920,
        height: 1080
    })

    await evaluate(page, function(wantedCategory, wantedTimeframe, slideData) {
        var origLoadNextSlide = window.apiExposer.prototype.loadNextSlide;
        window.apiExposer.prototype.loadNextSlide = function(){
            let e = this.sliderConfig.screens[this.screenId];
            let t = Object.values(this.sliderConfig.screens)[0];
            let slides = (e || t).slides;
            // pull an example slide object from https://media.ggleap.com/Centers/476fc5ba-6114-4445-94f8-b3734e7f770d/Slider/config.json
            console.log("Hello wolrld?")

            console.log("Grabbing from slideconfigs?")
            //let wantedSlideObj = slideConfigs.get("PUBG")

            /*
            for (var i = 0; i < slides.length; i++) {
                console.log("slide category:", slides[i].category);
                console.log("wantedCategory:", wantedCategory);

                if (slides[i].category == wantedCategory) {
                    wantedSlideObj = slides[i];
                    wantedSlideObj.timeframe = wantedTimeframe;
                    break
                }
            }
            */

            let realObj = t;
            if (e) {
                realObj = e;
            }

            realObj.slides = [slideData];

            console.log("in loadnextslide. sliderConfig: ", JSON.stringify(this.sliderConfig));
            origLoadNextSlide.apply(this, arguments);
        }
    }, wantedCategory, wantedTimeframe, getGameSlide(game));
    await timeout(2000);
    await page.screenshot({ path: filename});
    browser.close();
}

function getGameSlide(game) {
    let arr =  [{
        "type": "Dota",
        "category": "Dota",
        "previewPath": "dota_coins.jpg",
        "id": "88095e37-f463-4be2-8d3e-9d242a7ced7b",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    },
    {
        "type": "PUBG",
        "category": "PUBG",
        "previewPath": "pubg_last_day.jpg",
        "id": "70a56fdf-6baf-4ae8-ba7f-b3c3f0fc2349",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    },
    {
        "type": "TFT",
        "category": "TFT",
        "previewPath": "tft_last_day.jpg",
        "id": "bb40285e-f78a-45cd-b53c-61709e10d7a5",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    },
    {
        "type": "APEX",
        "category": "APEX",
        "previewPath": "apex_last_day.jpg",
        "id": "a54c5a3a-c083-4008-9d91-09cbed13e844",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    },
    {
        "type": "FortniteRankings",
        "category": "Fortnite",
        "previewPath": "fortnite_global.png",
        "id": "c1d92181-ac01-4db9-b15a-afc20b4fd4c9",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    },
    {
        "type": "Lol",
        "category": "LoL",
        "previewPath": "lol_coins.jpg",
        "id": "8bf48d4a-e29e-4c07-b554-8e977935b710",
        "mode": "static",
        "ranking": "regional",
        "displayValue": "coins",
        "duration": 30,
        "numberOfStaticRows": 10,
        "numberOfScrollingRows": 25,
        "timeframe": wantedTimeframe,
        "scrollRepeats": 10,
        "scrollSpeed": "normal",
        "backgroundUrl": null
    }]

    switch(game) {
        case "Dota":
            return arr[0]
        case "PUBG":
            return arr[1]
        case "TFT":
            return arr[2]
        case "APEX":
            return arr[3]
        case "Fornite":
            return arr[4]
        case "LOL":
            return arr[5]
    }
}
run();