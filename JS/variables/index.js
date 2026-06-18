// Hoisting
// Js moves declarations to top internally

//ex
console.log(a); 
var a = 10;

// JS see

var a;

console.log(a);
a = 10;

// let Hoisting

// console.log(b);
// let b = 10; // RefernceError



// Temporal Dead Zone (TDZ)

// {
//     console.log(age);

//     // TDZ starts here
//     let age = 22;
    
// }


// Scope Chain

const globalVar = "A";

function outer() {
    const outerVAr = "B";

    function inner() {
        console.log(globalVar);
        console.log(outerVAr);
    }

    inner()
}

outer();  // inner -> outer -> global



// THis keyword 

// this = who called the function

const user = {
    name: "Saurabh",

    greet() {
        console.log(`hello ${this.name}`);
    }
};

user.greet();

//  call()
// Invoke immerdiately

function greet() {
    console.log(this.name);
}

const user1 = {
    name: "Saurabh",
};

greet.call(user1);


// apply()

function intro(city, country) {
    console.log(
        this.name, // this is equal to name: saurabh
        city,
        country
    );
}

intro.apply({name: "Saurbah"}, ["Mumbai","India"]);


// bind()

function fuck() {
    console.log(`Fucke ${this.name}`);
}

const bound = fuck.bind({
    name: "saurabh"
}
);

bound();


// multi-line String

const msg = `
  hello
  world
  I'm here
`;

console.log(msg);


// length(), Index Acceess

const na = "Saurabh";
console.log(na.length);
console.log(na[1]);


