var url = "/api/user/login";
var params = {username: "admin" , password: "123456"};
var xhr = new XMLHttpRequest();
xhr.open("POST", url, true);
xhr.setRequestHeader("Content-Type", "application/json");
xhr.onload = function (e) {
  if (xhr.readyState === 4) {
    if (xhr.status === 200) {
      console.log(xhr.responseText);
    } else {
      console.error(xhr.statusText);
    }
  }
};
xhr.onerror = function (e) {
  console.error(xhr.statusText);
};
xhr.send(JSON.stringify(params));
// 登录

var url = "/api/recommends/7/comments";
var params = {userName: "admin", userId: "c4ca4238a0b923820dcc509a6f75849b", problemId: "2" , recommendReason: "Yes, lmj yyds"};
var xhr = new XMLHttpRequest();
xhr.open("POST", url, true);
xhr.setRequestHeader("Content-Type", "application/json");
xhr.setRequestHeader("Authorization", "orzzzzz");
xhr.onload = function (e) {
  if (xhr.readyState === 4) {
    if (xhr.status === 200) {
      console.log(xhr.responseText);
    } else {
      console.error(xhr.statusText);
    }
  }
};
xhr.onerror = function (e) {
  console.error(xhr.statusText);
};
xhr.send(JSON.stringify(params));
// 未登录

var url = "/api/recommends/23333/comments";
var params = {userName: "admin", userId: "c4ca4238a0b923820dcc509a6f75849b", problemId: "2" , recommendReason: "Yes, lmj yyds"};
var xhr = new XMLHttpRequest();
xhr.open("POST", url, true);
xhr.setRequestHeader("Content-Type", "application/json");
xhr.setRequestHeader("Authorization", "c81e728d9d4c2f636f067f89cc14862c");
xhr.onload = function (e) {
  if (xhr.readyState === 4) {
    if (xhr.status === 200) {
      console.log(xhr.responseText);
    } else {
      console.error(xhr.statusText);
    }
  }
};
xhr.onerror = function (e) {
  console.error(xhr.statusText);
};
xhr.send(JSON.stringify(params));
// 登录，推荐 id 错误

var url = "/api/recommends/7/comments";
var params = {userName: "admin", userId: "c4ca4238a0b923820dcc509a6f75849b", problemId: "2" , recommendReason: "Yes, lmj yyds"};
var xhr = new XMLHttpRequest();
xhr.open("POST", url, true);
xhr.setRequestHeader("Content-Type", "application/json");
xhr.setRequestHeader("Authorization", "c81e728d9d4c2f636f067f89cc14862c");
xhr.onload = function (e) {
  if (xhr.readyState === 4) {
    if (xhr.status === 200) {
      console.log(xhr.responseText);
    } else {
      console.error(xhr.statusText);
    }
  }
};
xhr.onerror = function (e) {
  console.error(xhr.statusText);
};
xhr.send(JSON.stringify(params));
// 登录，推荐 id 错误