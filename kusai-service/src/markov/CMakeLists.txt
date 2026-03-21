add_library(markov
  MarkovServiceImpl.cpp
)

target_include_directories(markov PUBLIC
  ${CMAKE_SOURCE_DIR}/include
)

target_sources(markov PUBLIC
  FILE_SET HEADERS
  BASE_DIRS ${CMAKE_SOURCE_DIR}/include
  FILES
  ${CMAKE_SOURCE_DIR}/include/kusai/service/markov/MarkovServiceImpl.hpp
)

target_link_libraries(markov
  PUBLIC
    proto
)
