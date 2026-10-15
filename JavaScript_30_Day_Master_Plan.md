# 🚀 JavaScript Master Plan — 30 Days, Zero to Confident

> **How to use in Notion:** Notion → Import → Text & Markdown → upload this file. Every `- [ ]` becomes a real checkbox. Tick them daily.
> **Daily routine (2–3 hrs):** 📖 Read theory (30 min) → 💻 Type every snippet yourself, never copy-paste (45 min) → 🛠 Practice tasks (60 min) → 📝 Write 3 lines "what I learned" (10 min).
> **Run code:** use the browser console (F12) or `node file.js`.

## 📊 Progress Tracker

| Week | Days | Theme | Status |
|---|---|---|---|
| 1 | 1–7 | Language foundations | ⬜ |
| 2 | 8–14 | Arrays, objects, functions | ⬜ |
| 3 | 15–21 | Classes, async, DOM | ⬜ |
| 4 | 22–30 | Events, fetch, modules, advanced, capstone | ⬜ |

---

# 🗓 WEEK 1 — Foundations

## Day 1 — Syntax, Strict Mode, Variables
**Theory:** A JS program is a list of *statements*. Semicolons end statements (optional but recommended; ASI inserts them and can surprise you). Comments are notes for humans. `"use strict"` turns silent mistakes into errors. `const` = value you won't reassign, `let` = value that changes, `var` = old style (avoid).
**Real world:** A shopping app stores `const TAX_RATE = 0.18` (never changes) and `let cartTotal = 0` (changes as you add items).

```js
"use strict";
// single-line comment
/* multi-line
   comment */
const TAX_RATE = 0.18;
let cartTotal = 0;
cartTotal = cartTotal + 500;
console.log(cartTotal * (1 + TAX_RATE)); // 590

// ASI trap: returns undefined!
function bad() {
  return
  { ok: true };
}
```

**Practice:**
- [ ] Declare your name, age, and a constant for your birth year; log a sentence.
- [ ] Make `x = 5` in strict mode with no declaration — read the error.
- [ ] Explain why `bad()` returns `undefined`.

**Checklist**
- [ ] Comments, statements, semicolons
- [ ] Strict mode `"use strict"`
- [ ] `var` vs `let` vs `const` (intro)

## Day 2 — Operators
**Theory:** Operators combine values. `==` converts types first (loose); `===` doesn't (strict) — **always prefer `===`**. `&&`/`||` return one of the operands, not just true/false. `??` falls back only for `null`/`undefined`; `||` falls back for any falsy value.
**Real world:** User settings: `volume = savedVolume ?? 50` keeps a saved `0` (muted) while `||` would wrongly reset it.

```js
// Arithmetic
console.log(10 + 3, 10 - 3, 10 * 3, 10 / 3, 10 % 3, 2 ** 5);
// Assignment
let score = 10; score += 5; score -= 2; score *= 2; score /= 2; score **= 2;
// Comparison
console.log(5 == "5", 5 === "5", 5 != "5", 5 !== "5", 5 > 3, 5 <= 5);
// Logical
const isAdult = true, hasId = false;
console.log(isAdult && hasId, isAdult || hasId, !isAdult);
// Bitwise (flags/permissions)
const READ = 1, WRITE = 2, EXEC = 4;
let perm = READ | WRITE;          // 3
console.log((perm & WRITE) !== 0); // true: has write
console.log(5 ^ 1, ~5, 1 << 3, 16 >> 2);
// Ternary
const age = 20;
console.log(age >= 18 ? "adult" : "minor");
// Nullish coalescing & optional chaining
const user = { profile: { city: "Varanasi" } };
console.log(user.address?.street ?? "No street");
console.log(0 || 50, 0 ?? 50); // 50, 0
// Spread / Rest
const nums = [1, 2, 3];
console.log(Math.max(...nums));
const sum = (...args) => args.reduce((a, b) => a + b, 0);
console.log(sum(1, 2, 3, 4));
```

**Practice:**
- [ ] Build a tip calculator using `*`, `%`, `+=`.
- [ ] Permission system with bitwise flags (read/write/delete).
- [ ] Safely read `order?.customer?.email ?? "unknown"`.

**Checklist**
- [ ] Arithmetic `+ - * / % **`
- [ ] Assignment `+= -=` etc.
- [ ] Comparison `== === != !== > < >= <=`
- [ ] Logical `&& || !`
- [ ] Bitwise operators
- [ ] Ternary `?:`
- [ ] Nullish coalescing `??`
- [ ] Optional chaining `?.`
- [ ] Spread / Rest `...`

## Day 3 — Control Flow & Error Handling
**Theory:** Control flow decides *which* code runs and *how often*. `for…of` loops over **values** of iterables (arrays, strings); `for…in` loops over **keys** of objects. `try/catch/finally` handles runtime errors; `finally` always runs.
**Real world:** Login: `if` checks password, `switch` picks a user role dashboard, `try/catch` handles bad JSON from a server.

```js
// if / else
const marks = 82;
if (marks >= 90) console.log("A");
else if (marks >= 75) console.log("B");
else console.log("C");

// switch
const role = "admin";
switch (role) {
  case "admin": console.log("Full access"); break;
  case "editor": console.log("Edit access"); break;
  default: console.log("Read only");
}

// loops
for (let i = 1; i <= 3; i++) console.log("for", i);
let n = 3; while (n > 0) { console.log("while", n); n--; }
let k = 0; do { console.log("do", k); k++; } while (k < 2);

for (const ch of "JS") console.log(ch);               // J, S
const car = { brand: "Tata", year: 2022 };
for (const key in car) console.log(key, car[key]);

// try / catch / finally
try {
  JSON.parse("{ bad json");
} catch (err) {
  console.log("Caught:", err.name);
} finally {
  console.log("Always runs");
}
```

**Practice:**
- [ ] FizzBuzz 1–100.
- [ ] Number guessing logic with `while`.
- [ ] Safe JSON parser that returns `null` on failure.

**Checklist**
- [ ] `if / else / switch`
- [ ] `for / while / do…while`
- [ ] `for…of / for…in`
- [ ] `try / catch / finally`

## Day 4 — Type System Basics
**Theory:** 7 primitives: `string, number, bigint, boolean, null, undefined, symbol`. Everything else is an object. **Falsy values:** `false, 0, -0, 0n, "", null, undefined, NaN`. Everything else is truthy (even `"0"`, `[]`, `{}`). JS converts types implicitly — know the explicit tools.
**Real world:** Form inputs are always strings: `"25" + 5` gives `"255"`. You must convert with `Number()`.

```js
console.log(typeof "hi", typeof 42, typeof 10n, typeof true,
            typeof undefined, typeof null /* "object" - famous bug */, typeof Symbol("id"));
console.log([] instanceof Array, new Date() instanceof Date);

// truthy / falsy
for (const v of [0, "", null, undefined, NaN, "0", [], {}]) {
  console.log(v, Boolean(v));
}

// conversions
console.log(Number("42"), Number("abc"), String(99), Boolean("hello"));
console.log(parseInt("42px"), parseInt("101", 2), parseFloat("3.14abc"));
console.log("25" + 5, "25" - 5, +"25" + 5);
```

**Practice:**
- [ ] Write `describe(value)` that returns type + truthiness.
- [ ] Convert a form's age string to number and validate.
- [ ] Predict then check: `[] + {}`, `"5" * "2"`, `null + 1`, `undefined + 1`.

**Checklist**
- [ ] Primitives (all 7)
- [ ] `typeof` / `instanceof`
- [ ] Truthy / falsy values
- [ ] `Number() String() Boolean() parseInt parseFloat`

## Day 5 — Variables & Scope, Hoisting, `this`
**Theory:** *Scope* = where a variable is visible. `var` is function-scoped; `let/const` are block-scoped. *Hoisting*: declarations are moved up; `var` is initialized as `undefined`, `let/const` sit in the **Temporal Dead Zone** until their line runs. `this` depends on **how** a function is called; `call/apply/bind` set it manually.
**Real world:** Borrowing a method: a `greet` function written for one user reused for another by `call`.

```js
console.log(a); var a = 1;          // undefined (hoisted)
// console.log(b); let b = 2;       // ReferenceError (TDZ)

{ var leaky = "I escape"; let safe = "I stay"; }
console.log(leaky);                 // works
// console.log(safe);               // error

const globalVar = "global";
function outer() {
  const outerVar = "outer";
  function inner() { console.log(globalVar, outerVar); }
  inner();
}
outer();

// this
const user = { name: "Asha", hi() { return "Hi " + this.name; } };
console.log(user.hi());
const loose = user.hi;              // lost this
// loose() -> undefined name (or error in strict mode)

function greet(greeting, punct) { return `${greeting}, ${this.name}${punct}`; }
const ravi = { name: "Ravi" };
console.log(greet.call(ravi, "Hello", "!"));
console.log(greet.apply(ravi, ["Hey", "?"]));
const boundGreet = greet.bind(ravi, "Namaste");
console.log(boundGreet("."));
```

**Practice:**
- [ ] Classic bug: `for (var i...) setTimeout` vs `let` — explain the difference.
- [ ] Fix a lost `this` using `bind`.
- [ ] Borrow `Array.prototype.slice` with `call` on `arguments`.

**Checklist**
- [ ] `var` vs `let` vs `const`
- [ ] Hoisting & Temporal Dead Zone
- [ ] Block, function, global scope
- [ ] `this`: `call`, `apply`, `bind`

## Day 6 — Strings
**Theory:** Strings are **immutable** — methods return new strings. Template literals use backticks, allow `${expr}` and multi-line text. Indexing starts at 0.
**Real world:** Generating usernames, slugs for URLs, masking card numbers, validating emails.

```js
const name = "  JavaScript Rocks  ";
console.log(`Hello ${name.trim()}! Length: ${name.trim().length}`);
console.log(name[2], name.at(-1));

const s = "JavaScript";
console.log(s.includes("Script"), s.startsWith("Java"), s.endsWith("pt"));
console.log(s.slice(0, 4), s.slice(-6), s.substring(4, 10), s.substr(4, 3)); // substr is deprecated
console.log(s.toUpperCase(), s.toLowerCase());
console.log("5".padStart(3, "0"), "Hi".padEnd(5, "."), "ab".repeat(3));

// split / replace / match
console.log("a,b,c".split(","));
console.log("2026-10-15".split("-").map(Number));
console.log("cat cat".replace("cat", "dog"), "cat cat".replaceAll("cat", "dog"));
console.log("Order #123 and #456".match(/\d+/g));

// slug generator
const slug = t => t.toLowerCase().trim().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
console.log(slug("  My First Blog Post! "));

// Unicode & escapes
console.log("Line1\nLine2", "Tab\tHere", "Quote: \"hi\"", "\u{1F600}", "😀".length);
// Mask a card
const card = "4111111111111234";
console.log(card.slice(-4).padStart(card.length, "*"));
// Regex basics: email check
console.log(/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test("me@site.com"));
```

**Practice:**
- [ ] Reverse a string and check palindromes.
- [ ] Count vowels in a sentence.
- [ ] Title-case a sentence ("hello world" → "Hello World").
- [ ] Mask a phone number: `98XXXXXX10`.

**Checklist**
- [ ] Template literals
- [ ] Length and index access
- [ ] `includes startsWith endsWith`
- [ ] `slice substring substr`
- [ ] `toUpperCase toLowerCase trim`
- [ ] `padStart padEnd repeat`
- [ ] `split replace replaceAll match`
- [ ] Unicode and escape sequences
- [ ] Regex basics

## Day 7 — Numbers, Math, BigInt + Week 1 Review
**Theory:** All numbers are 64-bit floats, so `0.1 + 0.2 !== 0.3`. Use `Number.isNaN` (not global `isNaN`). Integers beyond `2^53 - 1` need `BigInt`.
**Real world:** Money — store paise/cents as integers; use `toFixed` only for display. Random OTPs use `Math.random`.

```js
console.log(0.1 + 0.2, (0.1 + 0.2).toFixed(2), Math.abs(0.3 - (0.1 + 0.2)) < Number.EPSILON);
console.log(Number.isNaN(NaN), Number.isNaN("abc"), isNaN("abc"), Number.isFinite(1 / 0));
console.log(Math.round(2.5), Math.floor(-2.5), Math.ceil(2.1), Math.min(3, 1, 2), Math.max(3, 1, 2));
const otp = Math.floor(100000 + Math.random() * 900000);
console.log(otp);
const randInt = (min, max) => Math.floor(Math.random() * (max - min + 1)) + min;
console.log(randInt(1, 6));
console.log(1_000_000, 10 % 3, 2 ** 10, Number.MAX_SAFE_INTEGER);
const big = 9007199254740993n;
console.log(big + 2n, typeof big, BigInt(10) * 3n);
// isEven using %
const isEven = n => n % 2 === 0;
```

**Practice:**
- [ ] EMI calculator.
- [ ] Dice roller & OTP generator.
- [ ] Factorial of 25 using BigInt.

**Checklist**
- [ ] Float precision, `Number.isNaN`, `Number.isFinite`
- [ ] `Math`: round floor ceil random min max
- [ ] Numeric separators `1_000_000`
- [ ] BigInt
- [ ] `%` and `**`

**🏁 Week 1 Mini-Project:** *Unit converter + BMI calculator* (console). Use input validation, `switch`, `try/catch`, rounding.
- [ ] Project done · [ ] Week 1 review of weak topics

---

# 🗓 WEEK 2 — Arrays, Objects, Functions

## Day 8 — Arrays I: Create, Read, Iterate, Search, Modify
**Theory:** Arrays are ordered lists (objects underneath). `push/pop/shift/unshift/splice` **mutate**; `slice` doesn't.
**Real world:** A to-do list or shopping cart.

```js
const a = [1, 2, 3];
const b = Array.of(7, 8), c = Array.from("abc"), d = Array.from({ length: 3 }, (_, i) => i * 2);
console.log(a[0], a.length, a.at(-1)); a[1] = 20; a.length = 2;

[10, 20].forEach((v, i) => console.log(i, v));
for (const v of [5, 6]) console.log(v);

const fruits = ["apple", "banana", "mango", "banana"];
console.log(fruits.indexOf("banana"), fruits.lastIndexOf("banana"), fruits.includes("mango"));
const users = [{ id: 1, n: "A" }, { id: 2, n: "B" }];
console.log(users.find(u => u.id === 2), users.findIndex(u => u.id === 2));

const cart = ["pen"];
cart.push("book");      // add end
cart.unshift("bag");    // add start
cart.pop(); cart.shift();
const list = [1, 2, 3, 4, 5];
list.splice(1, 2, "x"); // remove 2 at index1, insert "x"
console.log(list, list.slice(1, 3));
```

**Practice:**
- [ ] Build add/remove/search for a to-do list.
- [ ] Remove duplicates using `indexOf`.
- [ ] Rotate an array by k positions.

**Checklist**
- [ ] Create: literals, `Array.of`, `Array.from`
- [ ] Read/update: index, length
- [ ] Iteration: `forEach`, `for…of`
- [ ] Search: `indexOf lastIndexOf includes find findIndex`
- [ ] Modify: `push pop shift unshift splice slice`

## Day 9 — Arrays II: Transform, Test, Sort, Destructure
**Theory:** `map` (same length, transformed), `filter` (subset), `reduce` (many → one). `sort` converts to strings by default — always pass a comparator for numbers.
**Real world:** Cart total = `reduce`; product listing = `filter` + `sort` + `map`.

```js
const products = [
  { name: "Phone", price: 20000, inStock: true },
  { name: "Cable", price: 300, inStock: true },
  { name: "TV", price: 45000, inStock: false },
];
const names = products.map(p => p.name);
const available = products.filter(p => p.inStock);
const total = products.reduce((sum, p) => sum + p.price, 0);
console.log(names, available.length, total);
console.log([[1, 2], [3, [4]]].flat(2), ["a b", "c"].flatMap(s => s.split(" ")));
console.log(products.some(p => !p.inStock), products.every(p => p.price > 100));

console.log([10, 9, 1].sort(), [10, 9, 1].sort((a, b) => a - b));
console.log([...products].sort((a, b) => b.price - a.price)[0].name);

const x = [1, 2], y = [3];
console.log(x.concat(y), [...x, ...y], [...x]); // copy
const [first, , third = 99, ...rest] = [1, 2, undefined, 4, 5];
console.log(first, third, rest);
let p = 1, q = 2; [p, q] = [q, p];   // swap
// group by (reduce)
const grouped = ["a", "bb", "cc"].reduce((acc, s) => ((acc[s.length] ??= []).push(s), acc), {});
```

**Practice:**
- [ ] Find the 2nd highest number.
- [ ] Group students by grade with `reduce`.
- [ ] Flatten nested categories and count items.

**Checklist**
- [ ] Transform: `map filter reduce flat flatMap`
- [ ] Test: `some every`
- [ ] Sort with comparators
- [ ] `concat`, spread, copy
- [ ] Array destructuring

## Day 10 — Objects I
**Theory:** Objects store key→value pairs. Access with `.` or `[]` (needed for dynamic keys). Destructuring unpacks properties.
**Real world:** A user profile, a product from an API.

```js
const key = "email";
const user = { name: "Asha", age: 25, [key]: "a@x.com", greet() { return `Hi ${this.name}`; } };
const name = "Ravi", age = 30;
const short = { name, age };                 // shorthand
console.log(user.name, user["age"], user[key], user.greet());
user.city = "Kashi"; delete user.age;

const { name: userName, city = "NA", phone = "none", ...others } = user;
console.log(userName, city, phone, others);
function show({ title, price = 0 }) { return `${title}: ₹${price}`; }
console.log(show({ title: "Book" }));

console.log(Object.keys(user), Object.values(user), Object.entries(user));
const merged = Object.assign({}, { a: 1 }, { b: 2 });
const copy = structuredClone({ deep: { x: 1 } });
console.log("name" in user, Object.hasOwn(user, "toString"));
// entries → object
console.log(Object.fromEntries(Object.entries({ a: 1, b: 2 }).map(([k, v]) => [k, v * 2])));
```

**Practice:**
- [ ] Count word frequency into an object.
- [ ] Deep-copy a nested settings object and change the copy.
- [ ] Convert array of `{id, name}` to a lookup `{[id]: name}`.

**Checklist**
- [ ] Literals & `.` / `[]` access
- [ ] Shorthand & computed keys
- [ ] Destructuring, defaults, rest
- [ ] `Object.keys/values/entries`
- [ ] `Object.assign`, `structuredClone`
- [ ] `hasOwn`, `in`

## Day 11 — Objects II: Prototypes, Descriptors, Freezing
**Theory:** Every object has a hidden link (`[[Prototype]]`) to another object; missing properties are looked up along this chain. Descriptors control `writable`, `enumerable`, `configurable`.
**Real world:** Config objects that must never change (`freeze`); inheriting shared behavior.

```js
const animal = { speak() { return `${this.name} makes a sound`; } };
const dog = Object.create(animal);
dog.name = "Tommy";
console.log(dog.speak(), Object.getPrototypeOf(dog) === animal);

const acct = {};
Object.defineProperty(acct, "id", { value: 101, writable: false, enumerable: false });
acct.id = 5; console.log(acct.id, Object.keys(acct));   // 101 []
console.log(Object.getOwnPropertyDescriptor(acct, "id"));

const config = Object.freeze({ api: "x", nested: { a: 1 } });
config.api = "y"; config.nested.a = 2;   // freeze is shallow
console.log(config);
const sealed = Object.seal({ a: 1 }); sealed.a = 2; sealed.b = 3; delete sealed.a;
const ne = Object.preventExtensions({ a: 1 }); ne.b = 2;
console.log(Object.isFrozen(config), Object.isSealed(sealed), Object.isExtensible(ne));
```

**Practice:**
- [ ] Build a prototype chain `vehicle → car → myCar`.
- [ ] Create a read-only constants object, try to break it.
- [ ] Write a deepFreeze function.

**Checklist**
- [ ] Prototypes & `Object.create`
- [ ] Property attributes & descriptors
- [ ] `freeze`, `seal`, `preventExtensions`

## Day 12 — Functions I
**Theory:** Declarations are hoisted fully; expressions are not. Arrow functions are short and **don't have their own `this`**. Default params fill missing values; rest collects extras.
**Real world:** Utility helpers: `formatPrice`, `calculateDiscount`, `validateEmail`.

```js
console.log(add(2, 3));                         // works: hoisted
function add(a, b) { return a + b; }
const sub = function (a, b) { return a - b; };  // expression
const mul = (a, b) => a * b;                    // arrow
const square = n => n * n;
const makeUser = (n, a) => ({ name: n, age: a }); // returning object literal

function price(amount, currency = "₹", ...taxes) {
  const t = taxes.reduce((s, x) => s + x, 0);
  return `${currency}${amount + t}`;
}
console.log(price(100), price(100, "$", 5, 3));

function divide(a, b) {
  if (b === 0) return "Cannot divide by zero";  // early return
  return a / b;
}
const timer = { s: 0, start() { setInterval(() => this.s++, 1000); } }; // arrow keeps this
```

**Practice:**
- [ ] Write `discount(price, pct = 10)`.
- [ ] Rewrite 3 functions in arrow form.
- [ ] Use early returns to flatten a nested validation.

**Checklist**
- [ ] Declarations vs expressions
- [ ] Arrow functions
- [ ] Default & rest params
- [ ] Return values & early returns

## Day 13 — Functions II: Closures, IIFE, Recursion, HOF
**Theory:** A **closure** is a function that remembers variables from where it was created. **HOF** = function taking/returning functions. Recursion needs a base case.
**Real world:** Closures power private counters, `debounce` for search boxes, and memoization caches.

```js
function counter() {
  let count = 0;                   // private
  return { inc: () => ++count, get: () => count };
}
const c = counter(); c.inc(); c.inc(); console.log(c.get()); // 2

(function () { const secret = "hidden"; console.log("IIFE ran"); })();

const fact = n => (n <= 1 ? 1 : n * fact(n - 1));
const fib = n => (n < 2 ? n : fib(n - 1) + fib(n - 2));
console.log(fact(5), fib(10));

// HOF + callbacks
const repeat = (n, fn) => { for (let i = 0; i < n; i++) fn(i); };
repeat(3, i => console.log("run", i));
const multiplier = f => n => n * f;
const double = multiplier(2);
console.log(double(8));

// debounce
function debounce(fn, ms) {
  let t; return (...a) => { clearTimeout(t); t = setTimeout(() => fn(...a), ms); };
}
// memoize
const memo = fn => { const cache = new Map(); return n => cache.has(n) ? cache.get(n) : (cache.set(n, fn(n)), cache.get(n)); };
```

**Practice:**
- [ ] Bank account with private balance via closure.
- [ ] Flatten nested arrays recursively.
- [ ] `once(fn)` that runs only the first time.

**Checklist**
- [ ] Closures & lexical environment
- [ ] IIFE
- [ ] Recursion
- [ ] Higher-order functions & callbacks

## Day 14 — Week 2 Review & Project
**🏁 Project:** *Expense Tracker (console)* — array of expense objects; functions to add, delete, filter by category, total with `reduce`, sort by amount, group by month.
- [ ] Project done
- [ ] Re-solve 5 weak exercises from Days 8–13
- [ ] Write your own 10-line cheat sheet of array methods

---

# 🗓 WEEK 3 — Classes, Async, DOM

## Day 15 — Classes I
**Theory:** `class` is cleaner syntax over prototypes. Methods live on the prototype (shared); fields set per instance. `#private` is truly inaccessible outside. `static` belongs to the class itself.
**Real world:** `BankAccount`, `Product`, `User` models.

```js
class BankAccount {
  #balance = 0;                    // private
  static bankName = "SBI";
  constructor(owner, initial = 0) { this.owner = owner; this.#balance = initial; }
  deposit(amt) { if (amt <= 0) throw new Error("Invalid"); this.#balance += amt; return this; }
  get balance() { return this.#balance; }
  set nickname(v) { this._nick = v.trim(); }
  static compare(a, b) { return a.balance - b.balance; }
}
const acc = new BankAccount("Asha", 1000).deposit(500);
console.log(acc.balance, BankAccount.bankName);
// acc.#balance -> SyntaxError
console.log(typeof BankAccount, Object.getPrototypeOf(acc) === BankAccount.prototype);
```

**Practice:**
- [ ] `Product` class with `price` getter/setter validation.
- [ ] `Counter` with private count and static `instances`.

**Checklist**
- [ ] `class`, `constructor`
- [ ] Instance vs prototype methods
- [ ] Getters/setters
- [ ] Static methods/fields
- [ ] Private fields `#name`

## Day 16 — Classes II: Inheritance & Composition
**Theory:** `extends` inherits; `super()` calls the parent. Override by redefining a method; `super.method()` reuses the parent's. **Prefer composition** ("has-a") over deep inheritance ("is-a"). Mixins add shared abilities.
**Real world:** `Employee → Manager`; a `CanFly` ability mixed into Bird and Plane.

```js
class Employee {
  constructor(name, salary) { this.name = name; this.salary = salary; }
  pay() { return this.salary; }
  toString() { return `${this.name}`; }
}
class Manager extends Employee {
  constructor(name, salary, bonus) { super(name, salary); this.bonus = bonus; }
  pay() { return super.pay() + this.bonus; }   // override
}
const m = new Manager("Ravi", 50000, 10000);
console.log(m.pay(), m instanceof Employee);   // chain: m → Manager → Employee → Object

// mixin
const CanFly = Base => class extends Base { fly() { return `${this.name} flies`; } };
class Bird { constructor(n) { this.name = n; } }
class Eagle extends CanFly(Bird) {}
console.log(new Eagle("Eddie").fly());

// composition
const logger = { log: m => console.log(`[LOG] ${m}`) };
class Service { constructor(logger) { this.logger = logger; } run() { this.logger.log("running"); } }
new Service(logger).run();
```

**Practice:**
- [ ] `Shape → Circle/Rectangle` with `area()` override.
- [ ] Make a `Serializable` mixin adding `toJSON()`.

**Checklist**
- [ ] `extends`, `super`
- [ ] Inheritance chains & overriding
- [ ] Mixins & composition basics

## Day 17 — Promises I
**Theory:** A Promise represents a future value: **pending → fulfilled** or **rejected** (settled once, forever). `then` handles success, `catch` failure, `finally` cleanup. Each `then` returns a new promise → chaining.
**Real world:** Ordering food: you get a "token" now (promise), the meal arrives later (resolve) or kitchen fails (reject).

```js
const wait = ms => new Promise(res => setTimeout(res, ms));
const orderFood = ok => new Promise((resolve, reject) => {
  setTimeout(() => ok ? resolve("🍕 Pizza") : reject(new Error("Kitchen closed")), 500);
});

orderFood(true)
  .then(food => { console.log("Got", food); return food + " + 🥤"; })
  .then(combo => console.log(combo))
  .catch(err => console.log("Error:", err.message))
  .finally(() => console.log("Done"));

Promise.resolve(42).then(console.log);
Promise.reject(new Error("no")).catch(e => console.log(e.message));

// error handling patterns: rethrow & recover
orderFood(false)
  .catch(e => { console.log("log:", e.message); return "🥪 Backup"; })
  .then(console.log);
```

**Practice:**
- [ ] Write `delay(ms)` and chain 3 delayed logs.
- [ ] Promisify a callback-style function.
- [ ] Make a promise that randomly fails; handle with catch + finally.

**Checklist**
- [ ] States & lifecycle
- [ ] `new Promise`
- [ ] `then / catch / finally`
- [ ] Error handling patterns
- [ ] `Promise.resolve / reject`

## Day 18 — async/await & Concurrency Helpers
**Theory:** `async` functions return promises; `await` pauses *that function* until the promise settles. Use `try/catch` for errors. Run independent tasks **in parallel** with `Promise.all`, not sequential awaits.
**Real world:** A dashboard loads profile, orders, and notifications simultaneously.

```js
const fetchUser = () => wait(300).then(() => ({ name: "Asha" }));
const fetchOrders = () => wait(500).then(() => [1, 2]);

async function load() {
  try {
    console.time("seq");
    const u = await fetchUser(); const o = await fetchOrders(); // 800ms
    console.timeEnd("seq");

    console.time("par");
    const [u2, o2] = await Promise.all([fetchUser(), fetchOrders()]); // 500ms
    console.timeEnd("par");
  } catch (e) { console.error(e); }
}
load();

const results = await Promise.allSettled([Promise.resolve(1), Promise.reject("x")]);
console.log(results);               // [{status:"fulfilled"...},{status:"rejected"...}]
console.log(await Promise.race([wait(100).then(() => "fast"), wait(500).then(() => "slow")]));
console.log(await Promise.any([Promise.reject(1), Promise.resolve("first success")]));
```
*(Top-level `await` works in modules/Node ESM; otherwise wrap in an async function.)*

**Practice:**
- [ ] Load 3 "APIs" in parallel and time it.
- [ ] Timeout wrapper with `Promise.race`.
- [ ] Retry function: try 3 times with delay.

**Checklist**
- [ ] `async / await`
- [ ] `Promise.all`
- [ ] `allSettled`, `race`, `any`

## Day 19 — Event Loop & AbortController
**Theory:** JS is single-threaded. The **call stack** runs sync code; **microtasks** (promises, `queueMicrotask`) run next; then one **macrotask** (timers, events). So promises beat `setTimeout(…, 0)`.
**Real world:** Why UI freezes during a heavy loop, and how cancelling a search request when the user types again saves bandwidth.

```js
console.log("1 sync");
setTimeout(() => console.log("4 macrotask"), 0);
Promise.resolve().then(() => console.log("3 microtask"));
console.log("2 sync");
// Output: 1, 2, 3, 4

// AbortController
const ctrl = new AbortController();
const task = new Promise((res, rej) => {
  const id = setTimeout(() => res("finished"), 2000);
  ctrl.signal.addEventListener("abort", () => { clearTimeout(id); rej(ctrl.signal.reason); });
});
setTimeout(() => ctrl.abort(new Error("Cancelled by user")), 500);
task.then(console.log).catch(e => console.log(e.message));
```

**Practice:**
- [ ] Predict output of 3 mixed sync/promise/timeout snippets, then run.
- [ ] Make a cancellable `sleep(ms, signal)`.

**Checklist**
- [ ] Microtasks vs macrotasks, event loop
- [ ] `AbortController` for cancellation

## Day 20 — DOM I: Tree, Selecting, Traversal
**Theory:** The browser turns HTML into a tree of **nodes**. Elements are one kind of node (text and comments are others). `querySelectorAll` returns a static `NodeList` (use `Array.from` for array methods); `getElementsBy*` return *live* collections.
**Real world:** Reading a form value, highlighting every "error" input, finding the card a clicked button belongs to.

Create `index.html` for practice:
```html
<ul id="list" class="menu"><li class="item">Home</li><li class="item active">About</li><li class="item">Contact</li></ul>
<script src="app.js"></script>
```
```js
const list = document.getElementById("list");
const items = document.getElementsByClassName("item");
const lis = document.getElementsByTagName("li");
const first = document.querySelector(".item");
const all = document.querySelectorAll(".item");
console.log(all instanceof NodeList, Array.isArray(all), Array.from(all).map(n => n.textContent));
console.log(first.closest("ul").id);

// traversal
console.log(list.children.length, list.firstElementChild.textContent);
console.log(first.nextElementSibling.textContent, first.parentNode === list);
// Node vs Element
console.log(list.childNodes.length, list.children.length); // text nodes counted only in childNodes
```

**Practice:**
- [ ] Log all link `href`s on any page via the console.
- [ ] Find the `.active` item's siblings.

**Checklist**
- [ ] Document tree: nodes, elements, attributes
- [ ] Node vs Element; `NodeList` vs Array
- [ ] Traversal: `parentNode children firstElementChild nextElementSibling`
- [ ] `getElementById / ClassName / TagName`
- [ ] `querySelector / querySelectorAll`
- [ ] `closest`

## Day 21 — DOM II: Create, Modify, Style, Measure (+ HTML Utilities)
**Theory:** Create with `createElement`, insert with `append/prepend/insertAdjacentHTML`. Use `textContent` for user text (safe); `innerHTML` with untrusted input causes **XSS**. Toggle styles via `classList`, not inline `style`, when possible.
**Real world:** Adding a task to a to-do list, dark mode toggle, showing a toast message.

```js
const li = document.createElement("li");
li.textContent = "New item";
li.className = "item";
li.classList.add("new"); li.classList.toggle("active"); li.classList.remove("new");
console.log(li.classList.contains("active"));
list.append(li);                 // end
list.prepend("Menu: ");          // start (string ok)
document.body.appendChild(document.createElement("hr"));
list.insertAdjacentHTML("beforeend", "<li class='item'>Blog</li>");
li.replaceWith(Object.assign(document.createElement("li"), { textContent: "Replaced" }));
list.lastElementChild.remove();

// attributes & properties
const input = document.createElement("input");
input.setAttribute("placeholder", "Name");
console.log(input.getAttribute("placeholder"), input.hasAttribute("disabled"));
input.removeAttribute("placeholder");
input.id = "name"; input.value = "Asha";
input.style.border = "2px solid red"; input.style.backgroundColor = "#eef";

// measurements
const box = list.getBoundingClientRect();
console.log(list.offsetWidth, list.clientHeight, list.scrollHeight, window.scrollY, box.top);
// safe vs unsafe
const evil = "<img src=x onerror=alert(1)>";
const p = document.createElement("p"); p.textContent = evil;   // safe
```

**Practice:**
- [ ] Build a to-do list that renders an array to the DOM.
- [ ] Dark-mode button toggling a class on `<body>`.
- [ ] Remove all items matching a condition.

**Checklist**
- [ ] `createElement`
- [ ] `append appendChild prepend`
- [ ] `insertAdjacentHTML`
- [ ] `replaceWith`, `remove`
- [ ] `get/set/has/removeAttribute`
- [ ] `id value textContent innerHTML`
- [ ] `className` and `classList`
- [ ] `style` inline styles
- [ ] Measurements: offset, client, scroll
- [ ] **HTML Element Utilities** (reading/writing, inserting, modifying, class utils, removing/replacing) — all covered above

---

# 🗓 WEEK 4 — Events, Fetch, Modules, Advanced, Capstone

## Day 22 — Events
**Theory:** Events bubble **up** from the target to ancestors (capturing goes down first). **Delegation:** attach one listener to a parent instead of many to children — works for dynamically added elements too.
**Real world:** One click handler on a product list that handles "Add to cart" for 1000 items.

```html
<form id="f"><input id="q"><button>Go</button></form>
<ul id="todos"><li>Task <button data-id="1">✖</button></li></ul>
```
```js
function onClick(e) {
  console.log(e.type, e.target, e.currentTarget, e.timeStamp);
}
btn.addEventListener("click", onClick);
btn.removeEventListener("click", onClick);

document.getElementById("f").addEventListener("submit", e => {
  e.preventDefault();                    // stop page reload
  console.log("Search:", document.getElementById("q").value);
});

// bubbling / capturing
outer.addEventListener("click", () => console.log("outer"), true);   // capture
inner.addEventListener("click", e => { e.stopPropagation(); console.log("inner"); });

// delegation
document.getElementById("todos").addEventListener("click", e => {
  const b = e.target.closest("button[data-id]");
  if (b) b.parentElement.remove();
});

q.addEventListener("input", e => console.log(e.target.value));
document.addEventListener("keydown", e => { if (e.key === "Escape") console.log("closed"); });
window.addEventListener("load", () => console.log("loaded"));
```

**Practice:**
- [ ] Live character counter on a textarea.
- [ ] Keyboard shortcuts (Ctrl+K opens search).
- [ ] To-do with delete via event delegation.

**Checklist**
- [ ] `addEventListener / removeEventListener`
- [ ] Event object: `target currentTarget type timeStamp`
- [ ] `preventDefault`
- [ ] Bubbling & capturing, `stopPropagation`
- [ ] Event delegation
- [ ] click, input, change, submit, keydown, load

## Day 23 — Fetch API I
**Theory:** `fetch` returns a promise that **only rejects on network failure** — HTTP 404/500 still resolve, so check `response.ok`. Parse the body with `json()`, `text()`, etc. (once only).
**Real world:** Loading posts from an API, submitting a signup form.

```js
const BASE = "https://jsonplaceholder.typicode.com";
async function getPosts() {
  const res = await fetch(`${BASE}/posts?_limit=3`);
  if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`);
  return res.json();
}
getPosts().then(console.log).catch(console.error);

const json = (url, method, body) => fetch(url, {
  method,
  headers: { "Content-Type": "application/json" },
  body: body ? JSON.stringify(body) : undefined,
}).then(r => r.ok ? r.json() : Promise.reject(new Error(r.status)));

json(`${BASE}/posts`, "POST", { title: "Hi", userId: 1 });     // POST
json(`${BASE}/posts/1`, "PUT", { id: 1, title: "New" });       // PUT
json(`${BASE}/posts/1`, "PATCH", { title: "Tweak" });          // PATCH
fetch(`${BASE}/posts/1`, { method: "DELETE" });                 // DELETE

// request init options
fetch(url, { method: "GET", headers: {}, mode: "cors", credentials: "same-origin", cache: "no-store" });
// other parsers
// await res.text(); await res.blob(); await res.arrayBuffer();
```

**Practice:**
- [ ] Render posts into the DOM with loading + error states.
- [ ] Create, update, delete a post (CRUD) with buttons.

**Checklist**
- [ ] GET, POST, PUT, PATCH, DELETE
- [ ] Request init: `method headers body mode credentials cache`
- [ ] `ok status statusText`
- [ ] `json() text() blob() arrayBuffer()`

## Day 24 — Fetch API II: Errors, JSON, URL, FormData, Streams, CORS, Abort
**Theory:** Retry only transient errors (network/5xx). `URLSearchParams` builds safe query strings. **CORS** is a browser rule: the *server* must allow your origin via headers; you can't fix it from frontend code.
**Real world:** Search box with cancel-on-type, file upload form, retrying flaky mobile networks.

```js
async function fetchRetry(url, tries = 3) {
  for (let i = 0; i < tries; i++) {
    try {
      const r = await fetch(url);
      if (r.ok) return r.json();
      if (r.status < 500) throw new Error(`Client error ${r.status}`); // don't retry 4xx
    } catch (e) { if (i === tries - 1) throw e; }
    await wait(2 ** i * 300);                 // exponential backoff
  }
}

console.log(JSON.stringify({ a: 1, d: new Date() }, null, 2), JSON.parse('{"x":[1,2]}'));

const u = new URL("https://shop.com/search?x=1");
u.searchParams.set("q", "red shoes"); u.searchParams.append("page", 2);
console.log(u.href, new URLSearchParams({ a: 1, b: 2 }).toString());

const fd = new FormData(document.querySelector("form"));
// fd.append("avatar", fileInput.files[0]); fetch("/upload", { method: "POST", body: fd });
console.log(Object.fromEntries(fd));

// stream basics
const res = await fetch(url);
const reader = res.body.getReader();
while (true) { const { done, value } = await reader.read(); if (done) break; console.log(value.length, "bytes"); }

// timeout via AbortController
const c = new AbortController();
const t = setTimeout(() => c.abort(), 5000);
try { await fetch(url, { signal: c.signal }); } catch (e) { if (e.name === "AbortError") console.log("timed out"); }
finally { clearTimeout(t); }
// AbortSignal.timeout(5000) is the shortcut
```

**Practice:**
- [ ] Debounced live search that aborts the previous request.
- [ ] Weather-style app using a free API with error UI.

**Checklist**
- [ ] Handling errors and retries
- [ ] `JSON.parse / stringify`
- [ ] `URL`, `URLSearchParams`
- [ ] `FormData`
- [ ] Streams (basics)
- [ ] CORS basics
- [ ] `AbortController` timeouts

## Day 25 — Modules
**Theory:** ES modules give each file its own scope. `export` what others may use; `import` what you need. Default export = one main thing; named exports = many. Bundlers (Vite, Rollup, Webpack) combine modules for production — concept only today.
**Real world:** Splitting an app into `api.js`, `utils.js`, `ui.js`.

`index.html`: `<script type="module" src="main.js"></script>` (serve via `npx vite` or Live Server)
```js
// utils.js
export const tax = 0.18;
export function format(n) { return `₹${n.toFixed(2)}`; }
export default class Cart { items = []; add(i) { this.items.push(i); } }

// main.js
import Cart, { tax, format as fmt } from "./utils.js";   // relative path
import * as utils from "./utils.js";
const cart = new Cart(); cart.add(100);
console.log(fmt(100 * (1 + tax)), utils.tax);

// dynamic import (lazy load)
button.addEventListener("click", async () => {
  const { default: Chart } = await import("./chart.js");
  new Chart();
});
// top-level await
const data = await fetch("/data.json").then(r => r.json());
```

**Concepts:** Vite = fast dev server + bundler; Rollup = library bundler; Webpack = highly configurable, older standard. They tree-shake unused exports and minify.

**Practice:**
- [ ] Split your Expense Tracker into 3 modules.
- [ ] Lazy-load a heavy module on button click.
- [ ] Run `npm create vite@latest` and inspect the output.

**Checklist**
- [ ] `export`, `export default`, `import`
- [ ] Module specifiers & paths
- [ ] Dynamic `import()`
- [ ] Top-level `await`
- [ ] Bundling overview (Vite, Rollup, Webpack)

## Day 26 — Timers, Date, Intl
**Theory:** Timers are macrotasks and not exact. `requestAnimationFrame` syncs with screen repaint (~60fps) — best for animation. `Date` months are **0-indexed**. `Intl` formats for any locale.
**Real world:** Countdown timer, "posted 3 hours ago", showing ₹ in Indian number format.

```js
const t = setTimeout(() => console.log("once"), 1000); clearTimeout(t);
let n = 0; const id = setInterval(() => { console.log(++n); if (n === 3) clearInterval(id); }, 500);
function animate(ts) { box.style.left = (ts / 10) % 300 + "px"; requestAnimationFrame(animate); }
requestAnimationFrame(animate);

const now = new Date();
const d = new Date(2026, 9, 15, 14, 30);          // Oct 15 (month 9!)
console.log(now.getFullYear(), now.getMonth(), now.getDay(), now.toISOString(), Date.now());
const diffDays = (new Date("2026-12-25") - now) / 86400000;

console.log(new Intl.DateTimeFormat("en-IN", { dateStyle: "full", timeStyle: "short" }).format(now));
console.log(new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR" }).format(1234567.5)); // ₹12,34,567.50
console.log(new Intl.NumberFormat("en", { notation: "compact" }).format(1500000)); // 1.5M
console.log(new Intl.RelativeTimeFormat("en", { numeric: "auto" }).format(-1, "day")); // yesterday
```

**Practice:**
- [ ] Countdown to New Year with `setInterval`.
- [ ] Age calculator from birth date.
- [ ] `timeAgo(date)` helper using `Intl.RelativeTimeFormat`.

**Checklist**
- [ ] `setTimeout`, `setInterval`, `requestAnimationFrame`
- [ ] `Date`, `Intl.DateTimeFormat`
- [ ] Intl number & date formatting

## Day 27 — Errors, Performance, Storage
**Theory:** Built-in errors: `Error, TypeError, ReferenceError, SyntaxError, RangeError`. Custom errors let you branch on `instanceof`. `localStorage` persists forever (strings only, ~5MB); `sessionStorage` dies with the tab; cookies go to the server with each request.
**Real world:** Remembering dark mode, keeping a cart after refresh, a custom `ValidationError` in a signup form.

```js
class ValidationError extends Error {
  constructor(field, msg) { super(msg); this.name = "ValidationError"; this.field = field; }
}
try { throw new ValidationError("email", "Invalid email"); }
catch (e) {
  if (e instanceof ValidationError) console.log(e.field, e.message);
  else throw e;                                    // don't swallow unknown errors
}
try { null.x; } catch (e) { console.log(e instanceof TypeError); }
try { new Error("outer", { cause: new Error("root") }); } catch {}

// performance
const t0 = performance.now();
for (let i = 0; i < 1e6; i++) Math.sqrt(i);
console.log(`${(performance.now() - t0).toFixed(2)} ms`);
console.time("x"); /* code */ console.timeEnd("x");

// storage
localStorage.setItem("theme", "dark");
console.log(localStorage.getItem("theme"));
localStorage.setItem("cart", JSON.stringify([{ id: 1 }]));
const cart = JSON.parse(localStorage.getItem("cart") ?? "[]");
sessionStorage.setItem("tab", "1"); localStorage.removeItem("theme"); localStorage.clear();
document.cookie = "visited=true; max-age=86400; path=/";
console.log(document.cookie);
```

**DevTools skills to try:** Console, Sources breakpoints, Network tab (see your fetch calls), Performance tab, Application tab (storage).

**Practice:**
- [ ] Persist the Day 21 to-do list in `localStorage`.
- [ ] Custom `HttpError` with status code.
- [ ] Find a slow function with `performance.now` and optimize it.

**Checklist**
- [ ] Error types & custom errors
- [ ] `performance.now`, DevTools
- [ ] `localStorage`, `sessionStorage`, cookies

## Day 28 — Web Workers & Regex Deep Dive
**Theory:** Heavy computation blocks the UI thread. A **Web Worker** runs JS on a separate thread and talks via `postMessage` (no DOM access). Regex: flags `g i m s u y`, groups, named groups, lookahead, lazy quantifiers.
**Real world:** Image processing / big-data parsing without freezing the page; password strength, date and phone validation.

```js
// worker.js
self.onmessage = e => {
  let sum = 0; for (let i = 0; i < e.data; i++) sum += i;
  self.postMessage(sum);
};
// main.js
const worker = new Worker("worker.js");
worker.postMessage(1e9);
worker.onmessage = e => console.log("Result:", e.data);
worker.onerror = e => console.error(e.message);
// worker.terminate();

// regex
const date = /(?<y>\d{4})-(?<m>\d{2})-(?<d>\d{2})/;
console.log("2026-10-15".match(date).groups.y);
console.log("a1b22c333".match(/\d+/g), [..."a1b22".matchAll(/\d+/g)].map(m => m.index));
console.log("price: $20".replace(/\$(\d+)/, "₹$1"));
console.log("John Smith".replace(/(\w+) (\w+)/, "$2, $1"));
console.log(/^(?=.*[A-Z])(?=.*\d).{8,}$/.test("Passw0rdX"));  // lookahead: strong password
console.log("<b>x</b><b>y</b>".match(/<b>.*?<\/b>/g));        // lazy
console.log(/^[6-9]\d{9}$/.test("9876543210"));               // Indian mobile
```

**Practice:**
- [ ] Offload a prime-number finder to a worker; keep the UI responsive.
- [ ] Regex validators: email, PIN code, hex color, URL.

**Checklist**
- [ ] Web Workers (intro)
- [ ] Regular expressions (deeper dive)

## Day 29 — Style, Linting, Testing
**Theory:** **Prettier** auto-formats; **ESLint** finds bugs/bad patterns. Tests prove code works and protect against regressions. **Vitest/Jest** pattern: Arrange → Act → Assert.
**Real world:** Teams run lint + tests on every commit so broken code never ships.

```bash
npm init -y
npm i -D prettier eslint vitest
npx prettier --write .
npx eslint --init
```
```js
// math.js
export const add = (a, b) => a + b;
export function divide(a, b) { if (b === 0) throw new Error("Div by zero"); return a / b; }

// math.test.js
import { describe, it, expect } from "vitest";   // Jest: globals, same API
import { add, divide } from "./math.js";
describe("math", () => {
  it("adds", () => expect(add(2, 3)).toBe(5));
  it("throws on zero", () => expect(() => divide(1, 0)).toThrow("Div by zero"));
  it("async", async () => expect(await Promise.resolve(1)).toBe(1));
});
// run: npx vitest
```
`package.json` scripts: `"test": "vitest", "lint": "eslint .", "format": "prettier --write ."`

**Practice:**
- [ ] Write 5 tests for your string utilities (slug, title-case).
- [ ] Fix 5 ESLint warnings in old code.

**Checklist**
- [ ] Prettier, ESLint (concepts + setup)
- [ ] Testing basics: Jest, Vitest

## Day 30 — 🏆 Capstone & Final Review
**Build:** *"Task & Weather Dashboard"* using **everything**:
- [ ] ES modules (`api.js`, `store.js`, `ui.js`, `main.js`)
- [ ] Classes (`Task`, `TaskManager`) with private fields & validation errors
- [ ] DOM rendering + event delegation (add, complete, delete tasks)
- [ ] `localStorage` persistence
- [ ] `fetch` weather/quote API with loading, error, retry, and abort
- [ ] `Intl` for dates and currency; `async/await` + `Promise.all`
- [ ] Debounced search/filter
- [ ] Prettier + ESLint + 5 Vitest tests
- [ ] Deploy free on Netlify/Vercel/GitHub Pages

**Final review:**
- [ ] Re-read all Week checklists; revisit unchecked items
- [ ] Explain closures, event loop, `this`, and prototypes out loud in 2 minutes each
- [ ] Solve 20 problems on LeetCode/Codewars (easy tier)

---

# 🎯 After 30 Days — Roadmap to "Master"

| Next Step | What to learn |
|---|---|
| TypeScript | Types for safer large apps |
| Framework | React (or Vue/Svelte) |
| Backend | Node.js, Express, REST, databases |
| Deep JS | "You Don't Know JS" (free), javascript.info |
| Practice | Build 1 project every 2 weeks, read others' code on GitHub |

> **Honest note:** 30 days gives you solid, job-ready fundamentals. *Mastery* comes from months of building real projects. Consistency beats intensity: 2 hours daily > 14 hours on Sunday.

**Useful free resources:** [javascript.info](https://javascript.info) · [MDN Web Docs](https://developer.mozilla.org) · [freeCodeCamp](https://freecodecamp.org)

## 📝 Daily Log Template (duplicate per day in Notion)
- **Date:** 
- **Day #:** 
- **Topics done:** 
- **Hardest concept:** 
- **3 things I learned:** 
- **Tomorrow's goal:** 
