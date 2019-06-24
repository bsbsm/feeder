var apiUrl = "http://127.0.0.1:8080/api/";

var offset = 0;
var pageSize = 20;

var newsList = document.getElementById('NewsList');
var searchBox = document.getElementById('SearchBox');

function showNews(newsCollection) {
    offset += newsCollection.length
    console.log("show news. New offset = ${offset}");
    removeNews()
 
    newsCollection.forEach(function(element) {
        var child = document.createElement('div');
        console.log(element);
        child.setAttribute("newsID", element.ID);
        child.setAttribute("class", "btn-link");
        child.onclick = function () {
            openNews(child);
        }
        child.innerHTML += element.Title;
        newsList.appendChild(child);
    });
}

function removeNews() {
    while (newsList.firstChild) {
        newsList.removeChild(newsList.firstChild);
    }
}

function openNews(newsItem) {
    var id = newsItem.getAttribute("newsID");

    console.log('details for news #' + id);

    var cb = function (responseText) {
        if (responseText == '' || responseText == null) {
            console.log("'details for news #" + id +"' response is null");
            return
        }

        console.log("'details for news #" + id +"' response: \n" + responseText);
        
        var news = JSON.parse(responseText)
        if (news != null) {
            document.getElementById("NewsListRow").style.display = "none";
        document.getElementById("NewsDetailRow").style.display='inline';
        document.getElementById("NewsDetailSpan").value = news.Title + " Source: " + news.Source; 
        document.getElementById("NewsDetailTArea").value = news.PayloadJSON;
        }
    };

    ajax.get(apiUrl + "news/" + id, null, cb, true);
}

function closeNews() {
    document.getElementById("NewsDetailSpan").innerHTML = "";
    document.getElementById("NewsDetailTArea").innerHTML = "";
    document.getElementById("NewsDetailRow").style.display = "none";
    document.getElementById("NewsListRow").style.display='inline';
}

//
// On load
//
window.onload = function () {

    if(newsList == null){
        newsList = document.getElementById('NewsList');
    }    

    if(searchBox == null){
        searchBox = document.getElementById('SearchBox');
    }    

        document.getElementById('SearchBtn').onclick = function () {
            console.log("'search news by title' request sending");

            var cb = function (responseText) {
                if (responseText == '' || responseText == null) {
                    console.log("'search news by title' response is null");
                    return
                }
    
                console.log("'search news by title' response: \n" + responseText);
                
                showNews(JSON.parse(responseText))
            };
            
            var title = searchBox.value; 
            // reset news offset
            offset = 0;
            var d = { t: title, off: offset, c: pageSize }
            
            ajax.get(apiUrl + 'news', d, cb, true);
        }

        document.getElementById('NextPageBtn').onclick = function () {
            console.log("'next news page' request sending");

            var cb = function (responseText) {
                if (responseText == '' || responseText == null) {
                    console.log("'next news page'  response is null");
                    return
                }
    
                console.log("'next news page'  response: \n" + responseText);
                showNews(JSON.parse(responseText))
            };

            var d = { off: offset, c: pageSize }
            
            ajax.get(apiUrl + 'news', d, cb, true);
        }

        document.getElementById('PrevPageBtn').onclick = function () {
            console.log("'previous news page' request sending");

            var cb = function (responseText) {
                if (responseText == '' || responseText == null) {
                    console.log("'previous news page'response is null");
                    return
                }
    
                console.log("'previous news page' response: \n" + responseText);
                showNews(JSON.parse(responseText))
            };

            offset -= pageSize
            
            if (offset < 0) {
                offset = 0
            }

            var d = {off: offset, c: pageSize }
            
            ajax.get(apiUrl + 'news', d, cb, true);
        }

        document.getElementById("AddFeedSourceBtn").onclick = function() {
            console.log("'add feed source' request sending");

            var cb = function (responseText) {
                if (responseText == '' || responseText == null) {
                    console.log("'add feed source' response is null");
                    return
                }
    
                console.log("'add feed source' response: \n" + responseText);
            };

            var sourceBox = document.getElementById("FeedSourceBox");
            var ruleBox = document.getElementById("RuleBox");

            var d = {u: sourceBox.value, r: ruleBox.value}
            
            ajax.put(apiUrl + 'feed', d, cb, true);
        }

        document.getElementById("FirstExample").onclick = function () {
            var cb = function(){}
            var d = {u: 'https://www.netroby.com/rss', r: 'Title=NewTitle'}
            
            ajax.put(apiUrl + 'feed', d, cb, true);
        }

        document.getElementById("SecondExample").onclick = function () {
            var cb = function(){}
            var d = {u: 'http://feeds.nytimes.com/nyt/rss/Technology', r: 'Title=title_field,Description=Body'}
            
            ajax.put(apiUrl + 'feed', d, cb, true);
        }

        document.getElementById("NewsDetailCloseBtn").onclick = closeNews;

        closeNews();
    }
    
    //
    // AJAX
    //
    var ajax = {};
    
    ajax.get = function (url, data, callback, async) {
        var query = [];
        for (var key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }

        ajax.send(url + (query.length ? '?' + query.join('&') : ''), callback, 'GET', null, async)
    };

    ajax.put = function (url, data, callback, async) {
        var query = [];
        for (var key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }

        ajax.send(url + (query.length ? '?' + query.join('&') : ''), callback, 'PUT', null, async)
    };

    ajax.send = function (url, callback, method, data, async) {
        console.log('sending request: "' + method + '" ' + url)

        if (async === undefined) {
            async = true;
        }
        var x = ajax.x();
        x.open(method, url, async);
        x.onreadystatechange = function () {
            if (x.readyState == XMLHttpRequest.DONE) { // XMLHttpRequest.DONE == 4
                callback(x.responseText);
            }
        };

        x.send(data);
    };

    ajax.x = function () {
        if (typeof XMLHttpRequest !== 'undefined') {
            return new XMLHttpRequest();
        }
        var versions = [
            "MSXML2.XmlHttp.6.0",
            "MSXML2.XmlHttp.5.0",
            "MSXML2.XmlHttp.4.0",
            "MSXML2.XmlHttp.3.0",
            "MSXML2.XmlHttp.2.0",
            "Microsoft.XmlHttp"
        ];
    
        var xhr;
        for (var i = 0; i < versions.length; i++) {
            try {
                xhr = new ActiveXObject(versions[i]);
                break;
            } catch (e) {
            }
        }
        return xhr;
    };