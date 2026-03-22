add_executable(kusai-service
  main.cpp
)

target_include_directories(kusai-service
  PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}
)

target_link_libraries(kusai-service
  PRIVATE
    proto
    graph
    markov
    registry
    textchain
    tokenizer
    gRPC::grpc++_reflection
)
