#include "kusai/service/textchain/TextChainServiceImpl.hpp"

namespace kusai::service {
grpc::Status TextChainServiceImpl::Create(grpc::ServerContext* context,
                                          const proto::textchain::v1::CreateRequest* request,
                                          proto::textchain::v1::CreateResponse* response) {
  const auto markov = markovRegistry_->get(request->markov_id());
  if (!markov) return {grpc::StatusCode::NOT_FOUND, "Markov not found"};

  const auto tokenizer = tokenizerRegistry_->get(request->tokenizer_id());
  if (!tokenizer) return {grpc::StatusCode::NOT_FOUND, "Tokenizer not found"};

  const auto chain = std::make_shared<TextChain>(*markov, *tokenizer);
  const auto id = textChainRegistry_->add(chain);
  response->set_textchain_id(id);

  return grpc::Status::OK;
}

grpc::Status TextChainServiceImpl::Delete(grpc::ServerContext* context,
                                          const proto::textchain::v1::DeleteRequest* request,
                                          proto::textchain::v1::DeleteResponse* response) {
  const auto success = textChainRegistry_->remove(request->textchain_id());
  response->set_success(success);

  return grpc::Status::OK;
}

grpc::Status TextChainServiceImpl::Train(grpc::ServerContext* context,
                                         const proto::textchain::v1::TrainRequest* request,
                                         proto::textchain::v1::TrainResponse* response) {
  const auto chain = textChainRegistry_->get(request->textchain_id());
  if (!chain) return {grpc::StatusCode::NOT_FOUND, "TextChain not found"};

  const std::vector seqs(request->sequences().begin(), request->sequences().end());
  (*chain)->train(seqs);
  response->set_success(true);

  return grpc::Status::OK;
}

grpc::Status TextChainServiceImpl::GenerateSequence(grpc::ServerContext* context,
                                                    const proto::textchain::v1::GenerateSequenceRequest* request,
                                                    proto::textchain::v1::GenerateSequenceResponse* response) {
  const auto chain = textChainRegistry_->get(request->textchain_id());
  if (!chain) return {grpc::StatusCode::NOT_FOUND, "TextChain not found"};

  const auto limit = request->has_limit() ? request->limit() : INT8_MAX;
  const auto sequence = (*chain)->generateSequence(request->context(), limit);
  response->mutable_tokens()->Assign(sequence.begin(), sequence.end());

  return grpc::Status::OK;
}

grpc::Status TextChainServiceImpl::GenerateText(grpc::ServerContext* context,
                                                const proto::textchain::v1::GenerateTextRequest* request,
                                                proto::textchain::v1::GenerateTextResponse* response) {
  const auto chain = textChainRegistry_->get(request->textchain_id());
  if (!chain) return {grpc::StatusCode::NOT_FOUND, "TextChain not found"};

  const auto limit = request->has_limit() ? request->limit() : INT8_MAX;
  const auto text = (*chain)->generateText(request->context(), limit);

  response->set_text(text);
  return grpc::Status::OK;
}
}  // namespace kusai::service
