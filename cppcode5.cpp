
#include <iostream>	
#include <iomanip>
#include <string>		
		
		template <typename T>
		void _format_output(std::ostream& out, const T& str) {
			// Pad the string with spaces so that it is at least 20 characters wide.
			out << str;
		}
		
					
				  
				#include <cstdio>



auto isPrime(int num) -> bool {
if (num <= 1) {
  return false;
 }

if (num == 2) {
  return true;
 }

if (num%2 == 0) {
  return false;
 }


auto maxDivisor = int(math.Sqrt(float64(num)));

for (auto i = 3; i <= maxDivisor; i += 2) {
if (num%i == 0) {
   return false;
  }

 }


 return true;
}


auto main() -> int {
 // Test for prime numbers
int numbers[] {2, 3, 17, 25, 29, 37}
;

for (auto  num  : numbers) {
if (isPrime(num)) {
printf("%d is prime.\n", num);
  } else {
printf("%d is not prime.\n", num);
  }

 }

return 0;
}
