let name = "saurabh";
console.log(name);
name = "me";
console.log(name);

let a = 10;
let b = 3;

console.log(a + b);
console.log(a - b);
console.log(a * b);
console.log(a / b);
console.log(a % b);
console.log(a + b);
console.log(a ** b);


let num = 8;
if (num % 2 == 0) {
    console.log("even");
} else{
    console.log("Odd");
}

// Assignment 

let count = 10;

count += 5;
count -= 2;
count *= 3;
count /= 2;
console.log(count);

let age = 20;
let citizen = true;
// Both must be true
if (age >= 18 && citizen) { 
    console.log("can vote ")  
}


// OR  One true is enough 
if (age >= 18 || citizen) {
    console.log("conditin passed");
}


// Ternary Operator

let status = age >= 18 ? "Adult" : "child";


//  Nullish Coalescing 
// Omnly replaces "null" && "undefined";
let score = 0;
console.log(score ?? 100);


// Optional Chaining ?.

let user = {};
// Give Error
// console.log(user.address.city);

// Undefined
console.log(user.address?.city);


// Spread Operator

// copy arrays
let nums  = [1, 2, 3];
let copy = [...nums];
console.log(copy)


// Merge arrays
let a1 = [1,2];
let b1 = [3,4];

let result = [...a1, ...b1];
console.log(result);


// Objects
let user1 = {
    name: "Saurabh"
};

let updated  = {
    ...user1,
    age: 22
};

console.log(updated);


//  Reat Operator ...

function sum(...numbers) {
    console.log(numbers);
}

sum(1,2,3,4);
