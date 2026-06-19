// Arrays

const fruits = ["apple", "banana", "guava"];
console.log(fruits[1]);


// Array.of()
// Creates array from arguments
// Array.of() always treat arguments as elements
const arr = Array.of(1,2,3,"saurah",5);

console.log(arr);


// Array() create lenght of array empty with that much size.
// but with mulitple values beahve like Array.
const arr1 = Array(5);
console.log(arr1);


// Array.from()
// Convert iterable into array

const a = Array.from({length:5}, (_, i) => i );
console.log(a);

const char = Array.from("hello");
console.log(char);

// for..of 
//best for values

for (const fruit of fruits){
    console.log(fruit);
}

//forEach()

fruits.forEach(
    fruit => {
        console.log(fruit)
    }
)


// diff between forEach and for..of
//forEach can't break for...of can break

for (const fruit of fruits ){
    if( fruit === "banana")
        break;
}



// indexOf()  lastIndexOf()

const arr2 = ["a", "b","c"];
console.log(arr2.indexOf("b"));

console.log(arr2.lastIndexOf("c"));

// includes()
console.log(arr2.includes("a"));

// find()

const users = [ {id:1, name: "Saurabh"}, {id:2, name:"raghuvanshi"} ];

const user = users.find(user => user.id === 2);

console.log(user);

// findIndex()

const index = users.findIndex(user => user.id === 2);
console.log(index);

