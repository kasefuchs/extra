#pragma once

#include "kusai/service/proto/tokenizer/v1.grpc.pb.h"

namespace kusai::service {
class TokenizerServiceImpl final : proto::tokenizer::v1::TokenizerService::Service {};
}  // namespace kusai::service
