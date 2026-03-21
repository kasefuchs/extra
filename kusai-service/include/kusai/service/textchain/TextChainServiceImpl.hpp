#pragma once

#include "kusai/service/proto/textchain/v1.grpc.pb.h"

namespace kusai::service {
class TextChainServiceImpl final : proto::textchain::v1::TextChainService::Service {};
}  // namespace kusai::service
