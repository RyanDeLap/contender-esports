var script = document.createElement('script');
// Javascript for Login Page

var formUser = document.getElementById('form_username');
var formPassword = document.getElementById('form_password');

var errormessage = document.getElementsByClassName('error_message');

function validUserAndPass(form) {
    data = {
        username: formUser.value.toLowerCase(),
        password: formPassword.value
    }

    $.post("/login", JSON.stringify(data), function (result, status) {
        errormessage[0].style.visibility = 'hidden';
        if (result.error == "") {
            window.location.href = "/dashboard";
        } else {
            errormessage[0].style.visibility = 'visible';
        }
    });
}
// Javascript for Home Page

// Javascript for the Hamburger Menu Change And Opening/Closing Menu
var styleElem = document.head.appendChild(document.createElement("style"));
var isMenuActive = false;
var container = document.getElementById('container2');

function openNav() {
    document.getElementById("sideNav").style.width = "400px";
}

function closeNav() {
    document.getElementById("sideNav").style.width = "0";
}

function changeHambugerMenu() {
    if (!isMenuActive) {
        openNav();
        styleElem.innerHTML = ".menu-btn:before {box-shadow: 0 0 0 #CCCCCC; background-color: #CCCCCC; transform: translateY(10px) rotate(45deg); } .menu-btn:after{ background-color: #CCCCCC; transform: translateY(-10px) rotate(-45deg);}";
        isMenuActive = true;
        document.body.classList.add("stop-scrolling");
    }
    else {
        closeNav();
        styleElem.innerHTML = ".menu-btn:before, .menu-btn:after { background-color: #424242; content: ''; display: block; height: 4px; transition: all 200ms ease-in-out; } .menu-btn:before { box-shadow: 0 10px 0 #424242; margin-bottom: 16px; }";
        isMenuActive = false;
        document.body.classList.remove("stop-scrolling");
    }
}

// Javascript for Closing Menu When You Click Off It
window.onload = function(){
    console.log("Hello World!")
    document.onclick = function(e){
        if(e.target.id !== 'sideNav' && e.target.id !== 'hambugerMenu'&& isMenuActive){
            //element clicked wasn't the Menu or Hambuger Icon; hide the menu
            changeHambugerMenu();
        }
    };
};

// Javascript for the Menu Data

//Variables
var webTitle = document.getElementById('title');
var menuItems = document.getElementById('menuItems');
var userForm = document.getElementById('userForm');
var userName = document.getElementById('form_userName');
var userPass = document.getElementById('form_userPassword');
var userPassCon = document.getElementById('form_userConfirmPassword');
var locationForm = document.getElementById('locationForm');
var styleElemLocation = document.head.appendChild(document.createElement("style"));
var locationTitleForm = document.getElementById('form_locationName');
var updateForm = document.getElementById('updateForm');
var titleUpdateForm = document.getElementById('form_titleUpdate');
var twitterForm = document.getElementById('twitterForm');
var facebookForm = document.getElementById('facebookForm');
var inputStringTwitter = document.getElementById('form_inputStringTwitter');
var inputStringFacebook = document.getElementById('form_inputStringFacebook');

// Everything to do with New User

function newUserShow() {
    userName.value = "";
    userPass.value = "";
    userPassCon.value = "";
    userForm.style.visibility = "visible";
    styleElemLocation.innerHTML = "#main {background-color: rgba(0, 0, 0, .4);}"
}

function addNewUser() {
    if (userPass.value === userPassCon.value && userPass.value !== "" && userPassCon.value !== "") {
        errormessage[0].style.visibility = 'hidden';
        userForm.style.visibility = "hidden";
        styleElemLocation.innerHTML = "#main {background-color: transparent;}"
        // TODO: Add new user to backend
    }
    else {
        errormessage[0].style.visibility = 'visible';
    }
}

function cancelNewUser() {
    userForm.style.visibility = "hidden";
    errormessage[0].style.visibility = 'hidden';
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
}

// Everything Doing with Location

var currentPage = null;

function Location(centerID, newTitle, newTwitter, newFacebook, newDiscordServer, centerData) {
    this.centerID = centerID;
    this.title = newTitle;
    this.twitter = newTwitter;
    this.facebook = newFacebook;
    this.discordServer = newDiscordServer;

    this.updatePage = () => {
        document.getElementById("config").style.display = 'block';
        currentPage = this;

        document.getElementById("twitter").innerHTML = ("<a class='twitter-timeline' data-width='500' data-height='900' href='" + this.twitter + "'>Twitter</a> <script async src='https://platform.twitter.com/widgets.js' charset='utf-8'></script>");
        document.getElementById("facebook").innerHTML = ("<div class='fb-page' data-href='" + this.facebook + "' data-tabs='timeline' data-width='500px' data-height='900px' data-small-header='false' data-adapt-container-width='true' data-hide-cover='false' data-show-facepile='true'><blockquote cite='" + this.facebook + "' class='fb-xfbml-parse-ignore'><a href='" + this.facebook + "'>Facebook</a></blockquote></div>");
        document.getElementById("discord").innerHTML = ("<iframe src='https://discordapp.com/widget?id=" + this.discordServer + "&theme=dark' width='500' height='900' allowtransparency='true' frameborder='0' id='discordiFrame'></iframe>");
        webTitle.innerHTML = (this.title);
        FB.XFBML.parse();
        twttr.widgets.load();
        this.reloadDiscordIFrame();

        console.log(document.getElementById("centerName").value )
        console.log(centerData.centerName)

        document.getElementById("centerName").value = centerData.centerName;
        document.getElementById("discordChannelID").value = centerData.discordChannelID;
        document.getElementById("discordGuildID").value = centerData.discordGuildID;
        document.getElementById("ggleapLink").value = centerData.ggLeapLink;

        //let supportedGames = ["Lol", "Dota", "APEX", "PUBG", "TFT", "Fortnite"]
        for(let i = 0; i < centerData.Schedules.length; i++) {
            let schedule = centerData.Schedules[i];
            let game = schedule.Game;

            console.log(schedule);

            document.getElementById("day" + game).value = schedule.DayOfWeek;
            document.getElementById("time" + game).value = schedule.TimeToPost;
        }
    }

    this.reloadDiscordIFrame = () => {
        var disc = document.getElementById("discordiFrame");
        if (disc) {
            disc.src="https://discordapp.com/widget?id=" + this.discordServer + "&theme=dark";
        }
    }
}

let supportedGames = ["LoL", "Dota", "APEX", "PUBG", "TFT", "Fortnite"]
function saveSchedules() {
    let centerID = this.currentPage.centerID;

    schedules = []
    for(let i = 0; i < supportedGames.length; i++) {
        let currentGame = supportedGames[i];
        let dayOfWeek = document.getElementById("day" + currentGame).value;
        if(dayOfWeek == null || dayOfWeek == "none"){
            continue;
        }
        let time = document.getElementById("time" + currentGame).value;
        schedules.push({game: currentGame, time: time, day: dayOfWeek})
    }

    console.log("Attempting to save: " + JSON.stringify(schedules))
    $.post("/save-schedules/"+centerID, JSON.stringify(schedules), function (result, status) {
        if (result.error == "") {
            alert("Saved schedules to backend.")
        } else {
            console.log("Result error" + result)
            alert("Error saving.")
        }
    });
}

function saveConfig() {
    let centerID = this.currentPage.centerID;
    let centerName = document.getElementById("centerName").value
    let discordChannelID = document.getElementById("discordChannelID").value
    let discordGuildID = document.getElementById("discordGuildID").value
    let ggLeapLink = document.getElementById("ggleapLink").value

    config = {centerName: centerName, discordChannelID: discordChannelID, discordGuildID: discordGuildID, ggLeapLink: ggLeapLink}

    $.post("/save-config/"+centerID, JSON.stringify(config), function (result, status) {
        if (result.error == "") {
            alert("Saved config.")
        } else {
            console.log("Result error" + result)
            alert("Error saving.")
        }
    });
}

var locations = [];

function printLocations() {
    var menuItems = document.getElementById('menuItems');
    menuItems.innerHTML = "";
    for (var i = 0; i < locations.length; i++) {
        menuItems.innerHTML += ("<a href='#' class='location' id=location'" + i + "'> " +  locations[i].title + " </a>");
    }

    var tempList = document.getElementsByClassName("location");

    for (var i = 0; i < tempList.length; i++) {
        tempList[i].addEventListener("click", locations[i].updatePage.bind(this, locations[i]));
    }
}

function newLocationShow() {
    locationForm.style.visibility = "visible";
    styleElemLocation.innerHTML = "#main {background-color: rgba(0, 0, 0, .4);}"
}

function addLocation() {
    let centerName = document.getElementById("newCenterName").value
    let ggLink = document.getElementById("newGGLeapLink").value
    let guildID = document.getElementById("newDiscordGuildID").value
    let channelID = document.getElementById("newDiscordChannelID").value

    data = {
        "centerName": centerName,
        "guildID": guildID,
        "channelID": channelID,
        "ggLink": ggLink
    }

    $.post("/create-center", JSON.stringify(data), function (result, status) {
        errormessage[0].style.visibility = 'hidden';
        if (result.error == "") {
            location.reload();
        } else {
            alert("Error adding center.")
            console.log(data)
        }
    });
}

function loadLocations() {
    locations = []
    console.log("Load locations?")
    $.get("/centers", function (result, status) {
        console.log(result)

        for (let i = 0; i < result.length; i++) {
            let center = result[i];
            console.log("Handling center", center)
            locations.push(new Location(center.centerID, center.centerName,
                "https://twitter.com/Contender_SGF",
                "https://www.facebook.com/contenderesports/",
                center.discordGuildID,
                center));
            printLocations();
        }
    });

    printLocations();
    locationForm.style.visibility = "hidden";
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
}

function cancelNewLocation() {
    locationForm.style.visibility = "hidden";
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
}

function updateLocationShow() {
    updateForm.style.visibility = "visible";
    styleElemLocation.innerHTML = "#main {background-color: rgba(0, 0, 0, .4);}"
    for (var i = 0; i < locations.length; i++) {
        if (webTitle.innerText === locations[i].title) {
            titleUpdateForm.value = locations[i].title;
        }
    }
}

function updateLocation() {
    updateForm.style.visibility = "hidden";
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
    for (var i = 0; i < locations.length; i++) {
        if (webTitle.innerText === locations[i].title) {
            locations[i].title = titleUpdateForm.value;
            locations[i].twitter = inputStringTwitter.value;
            locations[i].facebook = inputStringFacebook.value;
            locations[i].discordServer = titleUpdateForm.value + " Discord Feed";
            locations[i].discordBot =  titleUpdateForm.value + " Discord Bot";
            locations[i].placement =  i;
            locations[i].updatePage();
        }
    }
    printLocations();
}

function deleteLocation() {
    updateForm.style.visibility = "hidden";
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
    for (var i = 0; i < locations.length; i++) {
        if (webTitle.innerText === locations[i].title) {
            var j = i + 1;
            var k = i - 1;
            locations.splice(i, 1);
            if (i === 1) {
                locations[j].updatePage();
            }
            else {
                locations[k].updatePage()
            }
        }
    }
    printLocations();
}

function cancelUpdateLocation() {
    updateForm.style.visibility = "hidden";
    styleElemLocation.innerHTML = "#main {background-color: transparent;}"
}

function enterTwitterShow() {
    twitterForm.style.visibility = "visible";
    inputStringTwitter.value = "";
}

function addTwitterWebsite() {
    twitterForm.style.visibility = "hidden";
    var twitterLabel = document.getElementsByClassName("twitter_label");

    twitterLabel[0].innerHTML = inputStringTwitter.value;
    twitterLabel[1].innerHTML = inputStringTwitter.value;
}

function enterFacebookShow() {
    facebookForm.style.visibility = "visible";
    inputStringFacebook.value = "";
}

function addFacebookWebsite() {
    var facebookLabel = document.getElementsByClassName("facebook_label");

    facebookLabel[0].innerHTML = inputStringFacebook.value;
    facebookLabel[1].innerHTML = inputStringFacebook.value;

    facebookForm.style.visibility = "hidden";
}

function cancelWebsite() {
    facebookForm.style.visibility = "hidden";
    twitterForm.style.visibility = "hidden";
}
loadLocations();
printLocations();

$("#loginwithtwitter").click(function(){
    var currentCenter = currentPage.centerID;
    console.log("CurrentCenter: " + currentCenter)
    console.log(currentPage)

    var daForm = $("<form />").attr("action", "/login_twitter").attr("method", "POST").append(
        $("<input />").attr("type", "hidden").attr("name", "centerID").attr("value", currentCenter)
    )

    $(document.body).append(daForm)
    daForm.submit();
});

function makeFBPageOption(page) {
    console.log(page);
    return $('<option>')
        .val(page.id)
        .text(page.name)
        .attr('data-token', page.access_token);
}

function showFBPages() {
    var pages = [];
    var $select = $('#facebook_page_id');
    var extendedToken = window.userToken;
    if ($select) {
        $select.show();
        $("#facebookSaveButton").show();
        FB.api("/me/accounts?access_token="+extendedToken, function(response) {
            if (response.data.length > 0) {
                // remove existing or placeholder pages from list
                $select.empty();
            }

            for (var i = 0; i < response.data.length; i++) {
                pages.push(makeFBPageOption(response.data[i]))
            }
            $select.append(pages);
        });
    }
}

$("#facebookLoginButton").click(function(){
    FB.getLoginStatus(function(response) {
        if (response.status === 'connected') {
            FB.logout(function() {
                doFacebookInit();
            })
        } else {
            doFacebookInit();
        }
    });
});

function doFacebookInit() {
    console.log("in facebook init")
    FB.login(function(response) {
        // if auth succeeds, show list of pages to choose from
        if (response.status === 'connected') {
            // todo: error checking
            $.post("/upconvert_token", {'user_token': response.authResponse.accessToken}, function(data){
                if (!data.token) {
                    alert("some error upconverting user token")
                    return
                }
                console.log("upconverted:", data);
                window.userToken = data.token;
                showFBPages();
            });
        }
    }, {
        scope: 'manage_pages,publish_pages,pages_show_list'
    });
}

$("#facebookSaveButton").click(function(){
    // if we dont have a token, return
    if (window.userToken == "") {
        return;
    }

    var $select = $('#facebook_page_id');
    var option = $select.find('option:selected');

    // if there isnt any selected items, return
    if (1 > option.length) {
        return;
    }

    var facebookPageName = option.text();
    var facebookPageID = option.val();
    var facebookPageToken = option.attr('data-token');

    //upload stuff to server
    $.post("/save_facebook_page", {'centerID': currentPage.centerID, 'page_name': facebookPageName, 'page_id': facebookPageID, 'page_token': facebookPageToken}, function(){
        alert("Saved facebook page to backend.")
    });
});

    $("#send_post").click(function(){
        // todo: auth, maybe post, etc
        $.get("/test_post")
    });
