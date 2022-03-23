/*
The MIT License

Copyright 2022 Tsuzu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

#include <iostream>
#include <vector>
#include <array>
#include <string>
#include <fstream>
#include <bitset>
#include <sstream>
#include <map>
#include <optional>
#include <algorithm>
#include <cmath>

std::vector<std::string> words = {
#include "./words/words.csv"
#include "./words/candidates.csv"  
};

std::vector<std::string> candidates = {
#include "./words/candidates.csv"  
};

void loadWords() {
  std::ifstream wifs("./words/words.txt");

  for(std::string line; std::getline(wifs, line);) {
    if (!line.size())
      continue;

    words.push_back(line);
  }

  std::ifstream cifs("./words/candidates.txt");

  for(std::string line; std::getline(cifs, line);) {
    if (!line.size())
      continue;

    candidates.push_back(line);
  }
}

std::string reversed(const std::string& s) {
  return std::string(s.rbegin(), s.rend());
}

struct CharStatus {
  int orangeMin = 0;

  std::bitset<5> green;
  std::bitset<5> gray;

  std::string string() const {
    return std::to_string(orangeMin) + " " + reversed(green.to_string()) + " " + reversed(gray.to_string());
  }
};

bool operator < (const CharStatus& lhs, const CharStatus& rhs) {
  if (lhs.orangeMin < rhs.orangeMin) {
    return true;
  } else if (lhs.orangeMin > rhs.orangeMin) {
    return false;
  }

  const auto greenCompare = int(lhs.green.to_ulong()) - int(rhs.green.to_ulong());

  if (greenCompare < 0) {
    return true;
  } else if (greenCompare > 0) {
    return false;
  }

  return lhs.gray.to_ulong() < rhs.gray.to_ulong();
}

using Word = std::array<char, 5>;

Word toWord(const std::string& s) {
  Word w;
  for (int i = 0; i < 5; ++i) {
    w[i] = s[i];
  }
  return w;
}

std::string toString(const Word& s)  {
  return std::string(s.begin(), s.end());
}

std::ostream& operator <<(std::ostream& os, const Word& w) {
  os << toString(w);
  return os;
}

struct Status: public std::array<CharStatus, 26> {
  std::string string() const {
    std::stringstream ss;
    char c = 'a';
    for (auto& cs : *this) {
      ss << c << ": " << cs.string() << "\n";
      ++c;
    }
    return ss.str();
  }

  std::optional<Status> merge(const Status& st) const {
    Status ret;
    int greenCount = 0;
    for (int i = 0; i < 26; ++i) {
      auto& cs = (*this)[i];
      auto& stcs = st[i];
      auto& retcs = ret[i];

      retcs.gray = cs.gray | stcs.gray;
      retcs.green = cs.green | stcs.green;
      retcs.orangeMin = std::max(cs.orangeMin, stcs.orangeMin);

      if (retcs.green.count() > static_cast<std::size_t>(retcs.orangeMin)) {
        return std::nullopt;
      }
      if (retcs.gray.count() + static_cast<std::size_t>(retcs.orangeMin) > 5) {
        return std::nullopt;
      }
      if ((retcs.gray & retcs.green).count()) {
        return std::nullopt;
      }
      
      greenCount += retcs.green.count();
    }

    if (greenCount > 5) {
      return std::nullopt;
    }

    return ret;
  }

  bool match(const Word& w) const {
    for (int i = 0; i < 26; ++i) {
      auto& cs = (*this)[i];

      int counter = 0;
      for (int j = 0; j < 5; ++j) {
        if (w[j] == 'a' + i) {
          if (cs.gray[j]) {
            return false;
          }
          ++counter;
        } else {
          if (cs.green[j]) {
            return false;
          }
        }
      }

      if (counter < cs.orangeMin) {
        return false;
      }
    }

    return true;
  }
};

std::ostream& operator <<(std::ostream& os, const Status& st) {
  os << st.string();
  return os;
}

bool operator == (const Status& lhs, const Status& rhs) {
  return lhs.string() == rhs.string();
}

Status parseStatus(const std::string& s) {
  Status st;

  bool greenFlag = false, orangeFlag = false;
  int i = 0;
  for(auto c : s) {
    if (c == 'G') {
      greenFlag = true;
      continue;
    }
    if (c == 'O') {
      orangeFlag = true;
      continue;
    }

    if (greenFlag) {
      st[c - 'a'].green.set(i);
      ++st[c - 'a'].orangeMin;
      greenFlag = false;
    }
    else if (orangeFlag) {
      ++st[c - 'a'].orangeMin;
      st[c - 'a'].gray.set(i);
      orangeFlag = false;
    }
    else {
      st[c - 'a'].gray.set(i);
    }
    ++i;
  }

  for (auto& cs : st) {
    if (static_cast<size_t>(cs.orangeMin) == cs.green.count() && cs.gray.count() != 0) {
      for (int i = 0; i < 5; ++i) {
        if (!cs.green[i]) {
          cs.gray.set(i);
        }
      }
    }
  }

  return st;
}

Status getStatus(Word input, Word ans) {
  Status st;

  for (int i = 0; i < 5; ++i) {
    if (input[i] == ans[i]) {
      st[ans[i] - 'a'].green[i] = true;
      ++st[ans[i] - 'a'].orangeMin;
      input[i] = '-';
      ans[i] = '-';
    }
  }

  for (int i = 0; i < 5; ++i) {
    if (input[i] == '-') {
      continue;
    }
    auto& cs = st[input[i] - 'a'];

    bool found = false;
    for (int j = 0; j < 5; ++j) {
      if (input[i] == ans[j]) {
        ++cs.orangeMin;
        ans[j] = '-';
        found = true;
        break;
      }
    }
    cs.gray[i] = true;
    
    if (found) {
      continue;
    }

    if (static_cast<std::size_t>(cs.orangeMin) == cs.green.count()) {
      cs.gray = ~cs.green;
    }
  }

  return st;
}

std::map<Status, std::vector<Word>> getCandidates(const Status& currentStatus, const std::vector<Word> words, const Word& input) {
  std::map<Status, std::vector<Word>> candidates;

  for (auto& w : words) {
    auto st = getStatus(input, w).merge(currentStatus);

    if (!st) {
      // std::cout << w << " " << input << std::endl;
      continue;
    }

    candidates[*st].push_back(w);
  }

  return candidates;
}

int main() {
  // loadWords();

  std::vector<Word> words, candidates;

  std::transform(::words.begin(), ::words.end(), std::back_inserter(words), [](const std::string& s) {
    return toWord(s);
  });

  std::sort(words.begin(), words.end());

  std::transform(::candidates.begin(), ::candidates.end(), std::back_inserter(candidates), [](const std::string& s) {
    return toWord(s);
  });

  std::sort(candidates.begin(), candidates.end());

  Status status{};
  for (;;) {
    if (candidates.size() == 1) {
      std::cout << candidates[0] << std::endl;

      return 0;
    }

    if (candidates.size() == 2) {
      std::string result;

      std::cout << candidates[0] << std::endl;
      std::cin >> result;
      std::cout << candidates[1] << std::endl;

      return 0;
    }
    
    std::vector<
      std::tuple<
        double, // entropy
        Word,
        std::map<Status, std::vector<Word>>
      >
    > results;

    for (auto& input: words) {
      // std::cout << input << std::endl;
      const auto cands = getCandidates(status, candidates, input);

      double avgEnp = 0.;
      for (auto& [st, ws]: cands) {
        if (ws.empty()) {
          continue;
        }

        const auto p = ws.size() / static_cast<double>(candidates.size());

        avgEnp += - p * std::log2(p);
      }

      results.emplace_back(
        avgEnp,
        input,
        cands
      );
    }

    std::sort(results.begin(), results.end(), [](const auto& lhs, const auto& rhs) {
      return std::get<0>(lhs) > std::get<0>(rhs);
    });

    // std::cout << "1/" << candidates.size() << " " << std::get<0>(results[0]) << ": " << toString(std::get<1>(results[0])) << std::endl;

    std::cout << std::get<1>(results[0]) << std::endl;

    std::string result;
    std::cin >> result;

    const auto parsed = parseStatus(result);
    const auto merged = parsed.merge(status);

    if (!merged) {
      std::cerr << "error!" << std::endl;

      return 1;
    }

    // std::cout << "Current status:\n" << *merged << std::endl;

    status = *merged;
    candidates = std::get<2>(results[0])[*merged];
  }
}
