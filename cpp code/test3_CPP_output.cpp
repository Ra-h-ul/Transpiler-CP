#include <bits/stdc++.h>

template <typename T>
	void _format_output(std::ostream& out, const T& str) 
	{	
		out << str;
	}

// Program to print the first 5 natural numbers


auto main() -> int {

  // for loop terminates when i becomes 6
for (auto i = 1; i <= 5; i++) {
_format_output(std::cout, i);
std::cout << std::endl;
  }



return 0;
}
