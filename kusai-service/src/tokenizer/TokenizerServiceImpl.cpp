#include "kusai/service/tokenizer/TokenizerServiceImpl.hpp"

#include <kusai/tokenizer/SimpleTokenizer.hpp>

namespace kusai::service {
grpc::Status TokenizerServiceImpl::Create(grpc::ServerContext* context,
                                          const proto::tokenizer::v1::CreateRequest* request,
                                          proto::tokenizer::v1::CreateResponse* response) {
  std::shared_ptr<AbstractTokenizer> tokenizer;
  switch (request->impl_case()) {
    case proto::tokenizer::v1::CreateRequest::kSimple:
      tokenizer = std::make_shared<SimpleTokenizer>();
      break;
    default:
      return {grpc::StatusCode::INVALID_ARGUMENT, "Unknown implementation"};
  }

  const auto id = tokenizerRegistry_->add(tokenizer);
  response->set_tokenizer_id(id);

  return grpc::Status::OK;
}

grpc::Status TokenizerServiceImpl::Delete(grpc::ServerContext* context,
                                          const proto::tokenizer::v1::DeleteRequest* request,
                                          proto::tokenizer::v1::DeleteResponse* response) {
  const auto success = tokenizerRegistry_->remove(request->tokenizer_id());
  response->set_success(success);

  return grpc::Status::OK;
}

grpc::Status TokenizerServiceImpl::Encode(grpc::ServerContext* context,
                                          const proto::tokenizer::v1::EncodeRequest* request,
                                          proto::tokenizer::v1::EncodeResponse* response) {
  const auto tokenizer = tokenizerRegistry_->get(request->tokenizer_id());
  if (!tokenizer) return {grpc::StatusCode::NOT_FOUND, "Tokenizer not found"};

  for (const auto ids = (*tokenizer)->encode(request->text()); const auto id : ids) response->add_tokens(id);
  return grpc::Status::OK;
}

grpc::Status TokenizerServiceImpl::Decode(grpc::ServerContext* context,
                                          const proto::tokenizer::v1::DecodeRequest* request,
                                          proto::tokenizer::v1::DecodeResponse* response) {
  const auto tokenizer = tokenizerRegistry_->get(request->tokenizer_id());
  if (!tokenizer) return {grpc::StatusCode::NOT_FOUND, "Tokenizer not found"};

  const std::vector ids(request->tokens().begin(), request->tokens().end());
  response->set_text((*tokenizer)->decode(ids));
  return grpc::Status::OK;
}
}  // namespace kusai::service
