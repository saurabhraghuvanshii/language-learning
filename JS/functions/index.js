// // funnctions 

// // functionn declaration 

// function greet() {
//     console.log("heloo");
// }

// greet();

// // Functions Expresion
// // function stored inside a variable

// const greet1 = function() {
//     console.log("hello");
// }

// greet1();


// // NAmed Function Expression
// const greet2 = function sayHello() {
//    console.log("hello");
// }

// greet2();


// // Arrow functions

// const greet3 = () => {

// };

// // Single Expression
// const add = (a,b) => a + b;


// // single Parameter

// const square = (num) => num * num;



// // Rest Parameters
// // Collect remaining arguments

// function sum(...nums) {
//     console.log(nums)
// }

// sum(1,2,3);

// // Return Values
// // Functions can return anything

// function add1(a, b) {
//     return a + b;
// }


// // Return object

// function createUser() {
//     return {
//         name: "Saurabh"
//     };
// }



// // Return functions
// function outer() {

// }


// for (var i = 0; i < 3; i++) {
//     setTimeout(() => console.log(i), 1);
// }

// for (let i = 0; i < 3; i++) {
//     setTimeout(() => console.log(i), 1);
// }


// let greeting;
// greeting = {};
// console.log(greeting);

// function bark() {
//     console.log('Woof!');
// }
  
// bark.animal = 'dog';


// function person (firstName, lastName) {
//     this.firstName = firstName;
//     this.lastName = lastName;
// }

// const member = new person("saurabh", "singh");

// person.prototype.getFullName = function() {
//     return `${this.firstName} ${this.lastName}`
// };


// console.log(member.getFullName());

// function Person(firstName, lastName) {
//     this.firstName = firstName;
//     this.lastName = lastName;
//   }
  
//   const lydia = new Person('Lydia', 'Hallie');
//   const sarah = Person('Sarah', 'Smith');
  
//   console.log(lydia);
//   console.log(sarah);

let number = 0;
console.log(number++);
console.log(++number);
console.log(number);