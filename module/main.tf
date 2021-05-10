variable "pet_count" {
  type    = number
  default = 1
}

resource "random_pet" "pets" {
  count = var.pet_count
}

output "pets" {
  value = random_pet.pets
}

output "message" {
  value = "this is message 1.0"
}
