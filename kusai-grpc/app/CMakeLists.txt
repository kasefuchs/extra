add_executable(kusai-grpc
  main.cpp
)

target_include_directories(kusai-grpc
  PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}
)

target_link_libraries(kusai-grpc
  PRIVATE
    proto
    registry
    kusai::textchain
)
