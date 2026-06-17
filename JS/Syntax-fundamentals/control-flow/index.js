// If else

let marks = 85;

if (marks >= 90) {
    console.log("A");
}else if (marks >= 80) {
    console.log("B");
}else {
    console.log("C");
}

// Switch

let day = 2;

switch (day) {
    case 1:
        console.log("Monday");
        break;
    case 2: 
        console.log("Tuesday");
        break;
    default:
        console.log("fuck");
}

// for loop
for (let i = 1; i <= 5; i++){
    console.log(i);
}


// while

let j = 1;

while(j < 5 ){
   console.log(j)
   j++;
}

// do while  Runs at least once

do {
    console.log(j);
}while (j < 5);


// for .. of

let fruits = ["apple", "banana"];

for (const fruit of fruits) {
    console.log(fruit);
}



// for  in

let user = {
    name: "saurabh",
    age: "21"
};

for (const key in user) {
    console.log(key);
    console.log(user[key]);
}


// try catch finally

try  {
    Json.parse("{");
} catch (error) {
    console.log("Invalid JSON");
} finally {
    console.log("Always runs");
}



// Strict Mode

// "use strict";

// // Error
// name = "saurabh";

const aa = typeof null;
console.log(aa);


// instanceof

let arr = [];
console.log(arr instanceof Array);



