#pragma once

#include <absl/container/flat_hash_map.h>

#include <memory>

namespace kusai::grpc {
template <typename T>
class Registry {
 public:
  using ValuePtr = std::shared_ptr<T>;

  Registry() = default;

  ~Registry() = default;

  [[nodiscard]] std::uint32_t add(const ValuePtr& obj) {
    auto id = nextId_++;
    items_.emplace(id, obj);
    return id;
  };

  [[nodiscard]] std::optional<ValuePtr> get(std::uint32_t id) const {
    if (const auto it = items_.find(id); it != items_.end()) {
      return it->second;
    }

    return std::nullopt;
  };

  [[nodiscard]] bool remove(std::uint32_t id) { return items_.erase(id); };

 private:
  std::uint32_t nextId_{0};
  absl::flat_hash_map<std::uint32_t, ValuePtr> items_;
};
}  // namespace kusai::grpc
