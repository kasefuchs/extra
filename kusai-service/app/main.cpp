#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>

#include <kusai/markov/AbstractMarkov.hpp>
#include <kusai/textchain/TextChain.hpp>
#include <kusai/tokenizer/AbstractTokenizer.hpp>
#include <memory>

#include "kusai/service/graph/GraphServiceImpl.hpp"
#include "kusai/service/markov/MarkovServiceImpl.hpp"
#include "kusai/service/textchain/TextChainServiceImpl.hpp"
#include "kusai/service/tokenizer/TokenizerServiceImpl.hpp"

int main() {
  const auto graphRegistry = std::make_shared<kusai::service::Registry<kusai::AbstractGraph>>();
  const auto tokenizerRegistry = std::make_shared<kusai::service::Registry<kusai::AbstractTokenizer>>();
  const auto markovRegistry = std::make_shared<kusai::service::Registry<kusai::AbstractMarkov>>();
  const auto textChainRegistry = std::make_shared<kusai::service::Registry<kusai::TextChain>>();

  kusai::service::GraphServiceImpl graphService(graphRegistry);
  kusai::service::TokenizerServiceImpl tokenizerService(tokenizerRegistry);
  kusai::service::MarkovServiceImpl markovService(graphRegistry, markovRegistry);
  kusai::service::TextChainServiceImpl textChainService(markovRegistry, tokenizerRegistry, textChainRegistry);

  grpc::reflection::InitProtoReflectionServerBuilderPlugin();

  grpc::ServerBuilder builder;
  builder.AddListeningPort("0.0.0.0:50051", grpc::InsecureServerCredentials());
  builder.RegisterService(&graphService);
  builder.RegisterService(&tokenizerService);
  builder.RegisterService(&markovService);
  builder.RegisterService(&textChainService);

  const std::unique_ptr server(builder.BuildAndStart());

  server->Wait();
}
