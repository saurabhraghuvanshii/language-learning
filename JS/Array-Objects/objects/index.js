// Objects
// objects  are key-value pairs

const user = {
    name: "saurabh",
    age: 21
};

// Access Properties
// Dot notation.
console.log(user.name);

// Bracket notation
console.log(user["name"]);

// Dyanmic keys 

const key = "name";
console.log(user[key]);



// Shorthand Properties

const name = "shivam";

const user1 = {
    name
};

console.log(user1.name);

// Computed Keys

const key1 = "email";

const user2 = {
    [key1]: "tes@gmai.com"
};


// Object Destructuring 

const user3 = {
    name1: "Shubham",
    age: 21
};

const {name1, age} = user3;

const { name1: UserName} = user3;

console.log(name1, UserName);


// Prototypes
// every objet has a prototype.

const use = {};

// Object.create()

const person = {
    greet() {
        console.log("Hello");
    }
};

const user4 =  Object.create(person);

user4.greet();


// Object.keys()

console.log(Object.keys(user3));

// Object.values()

console.log(Object.values(user3));


// Object.entries()

console.log(Object.entries(user3));

// Onject.assign()
// Merge objects

// const Merged = Object.assign({}, obj1, obj2);

// modern 

// const merge = { ...obj1, ...obj2};


// structuredClone
// Deeep clone

const copy = structuredClone(user3);

// hasOwn()
// check own property

Object.hasOwn( user3, "name1");



// in Operator 
// check prototypes chain too 

console.log("name1" in user3);

// freeze()
// Completely immutable
Object.freeze(user3);



// seal()
// Can update exiting.
//Cannot add/remove

Object.seal(user3);
//

// preventExtensions()
//cannot addd new properties 
// can modify existing.

Object.preventExtensions(
    user1
);

