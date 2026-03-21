add_library(tokenizer
  TokenizerServiceImpl.cpp
)

target_include_directories(tokenizer PUBLIC
  ${CMAKE_SOURCE_DIR}/include
)

target_sources(tokenizer PUBLIC
  FILE_SET HEADERS
  BASE_DIRS ${CMAKE_SOURCE_DIR}/include
  FILES
  ${CMAKE_SOURCE_DIR}/include/kusai/service/tokenizer/TokenizerServiceImpl.hpp
)

target_link_libraries(tokenizer
  PUBLIC
    proto
)
