wordle: wordle.cpp
	g++ -O3 -o wordle -std=c++17 wordle.cpp

.PHONY: download
download:
	curl -o ./words/candidates.csv ${DOWNLOAD_CANDIDATES_URL}
	curl -o ./words/words.csv ${DOWNLOAD_WORDS_URL}
