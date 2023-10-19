#include <bits/stdc++.h>

template <typename T>
void _format_output(std::ostream &out, const T &str)
{
	out << str;
}

auto main() -> int
{
	auto x = 10;
	auto y = 6;
	auto z = x + y;
	_format_output(std::cout, z);
	std::cout << std::endl;

	return 0;
}
