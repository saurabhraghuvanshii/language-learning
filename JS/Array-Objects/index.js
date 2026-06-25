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




// transform methods

// map()

const nums = [1,2,3];

const doubled = nums.map(num => num * 2);

console.log(doubled);


// filter
// keep matching elements

const nums1 = [1,2,3,4];

const even = nums1.filter(num => num % 2 == 0 );

console.log(even);


// reduce()
// most powerful

const nums2 = [1,2,3,4];
const sum = nums.reduce((acc, curr) => acc + curr, 0);

console.log(sum);



// flat()

const arr4 = [[1,2], [3,4]];

console.log(arr4.flat());

// flatMap() 
// Map + flat 

const words = ["hello world", "I'm boss"];

const result = words.flatMap((word) => word.split(" "));

console.log(result);


// Testing Arrays

// some()
// At least one

const nums3 = [1,2,3];

console.log(nums3.some(num => num > 2));


// every()
// All must pass.
console.log(nums3.every( num => num > 0));



// Modify Arrays

// push()
// Add end.

const arr5 = [1,2];
arr5.push(3);

console.log(arr5);


// pop()
// Remove end.
console.log(arr5.pop());

// shift()
// Remove first
console.log(arr5.shift());


// unshift()
// Add first. and returen length of array 
const arr6 = [1,2,3,4,5]
console.log(arr6.unshift(0));

// splice()
// change original array
const arr7 = [1,2,3,4,5,6,7];
// remove
console.log(arr7.splice(1,2));

// Insert

console.log(arr7.splice(1,0,"new"));


// slice()
// Does not modify original

console.log(arr7.slice(1,3));



// Sort

[10,2,5].sort();
// ans [10,2,5] because it convert it into strings rather than a number lexicographically

const sr = [10,2,5].sort((a,b) => a-b);
// descending b-a;
console.log(sr);


// Concat

const n = [1,3,4];
const n1 = [2,4];

console.log(n.concat(n1));

// Spread
const merged = [ ...n, ...n1];
console.log(merged);


// Array Destructing

const colors = ["red", "green", "blue"];

const [first, second] = colors;
console.log(first,second);