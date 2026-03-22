#pragma once

#include <kusai/tokenizer/AbstractTokenizer.hpp>

#include "kusai/service/proto/tokenizer/v1.grpc.pb.h"
#include "kusai/service/registry/Registry.hpp"

namespace kusai::service {
class TokenizerServiceImpl final : public proto::tokenizer::v1::TokenizerService::Service {
 public:
  explicit TokenizerServiceImpl(const std::shared_ptr<Registry<AbstractTokenizer>>& tokenizerRegistry)
      : tokenizerRegistry_(tokenizerRegistry) {}

  ~TokenizerServiceImpl() override = default;

  grpc::Status Create(grpc::ServerContext* context, const proto::tokenizer::v1::CreateRequest* request,
                      proto::tokenizer::v1::CreateResponse* response) override;

  grpc::Status Delete(grpc::ServerContext* context, const proto::tokenizer::v1::DeleteRequest* request,
                      proto::tokenizer::v1::DeleteResponse* response) override;

  grpc::Status Encode(grpc::ServerContext* context, const proto::tokenizer::v1::EncodeRequest* request,
                      proto::tokenizer::v1::EncodeResponse* response) override;

  grpc::Status Decode(grpc::ServerContext* context, const proto::tokenizer::v1::DecodeRequest* request,
                      proto::tokenizer::v1::DecodeResponse* response) override;

 private:
  std::shared_ptr<Registry<AbstractTokenizer>> tokenizerRegistry_;
};
}  // namespace kusai::service
