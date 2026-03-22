#pragma once

#include <kusai/textchain/TextChain.hpp>

#include "kusai/service/proto/textchain/v1.grpc.pb.h"
#include "kusai/service/registry/Registry.hpp"

namespace kusai::service {
class TextChainServiceImpl final : public proto::textchain::v1::TextChainService::Service {
 public:
  explicit TextChainServiceImpl(const std::shared_ptr<Registry<AbstractMarkov>>& markovRegistry,
                                const std::shared_ptr<Registry<AbstractTokenizer>>& tokenizerRegistry,
                                const std::shared_ptr<Registry<TextChain>>& textChainRegistry)
      : markovRegistry_(markovRegistry), tokenizerRegistry_(tokenizerRegistry), textChainRegistry_(textChainRegistry) {}

  ~TextChainServiceImpl() override = default;

  grpc::Status Create(grpc::ServerContext* context, const proto::textchain::v1::CreateRequest* request,
                      proto::textchain::v1::CreateResponse* response) override;

  grpc::Status Delete(grpc::ServerContext* context, const proto::textchain::v1::DeleteRequest* request,
                      proto::textchain::v1::DeleteResponse* response) override;

  grpc::Status Train(grpc::ServerContext* context, const proto::textchain::v1::TrainRequest* request,
                     proto::textchain::v1::TrainResponse* response) override;

  grpc::Status GenerateSequence(grpc::ServerContext* context,
                                const proto::textchain::v1::GenerateSequenceRequest* request,
                                proto::textchain::v1::GenerateSequenceResponse* response) override;

  grpc::Status GenerateText(grpc::ServerContext* context, const proto::textchain::v1::GenerateTextRequest* request,
                            proto::textchain::v1::GenerateTextResponse* response) override;

 private:
  std::shared_ptr<Registry<AbstractMarkov>> markovRegistry_;
  std::shared_ptr<Registry<AbstractTokenizer>> tokenizerRegistry_;
  std::shared_ptr<Registry<TextChain>> textChainRegistry_;
};
}  // namespace kusai::service
