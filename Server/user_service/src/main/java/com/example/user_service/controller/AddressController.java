//package com.example.user_service.controller;
//
//
//import org.springframework.web.bind.annotation.*;
//
//import java.util.List;
//import java.util.Map;
//
//
//@RestController
//@RequestMapping("/a")
//public class AddressController {
//    @GetMapping("/hello")
//    public String hello() {
//        System.out.println("Nhan A");
//        return "Hello from Service A";
//    }
//}
////
////
////    @Autowired
////    private UserService userService;
////
////    @Autowired
////    private AddressService addressService;
////
////    @Autowired
////    private CustomerService customerService;
////
////    @GetMapping()
////    public ResponseEntity<?> getUserInfo(@RequestHeader("Authorization") String token) {
////        try {
////            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
////            Users user = userService.getUserById(userId);
////            Customer customer = customerService.getCustomerById(userId);
////            if (user == null) {
////                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
////                        "code", 404,
////                        "message", "Không tìm thấy người dùng"
////                ));
////            }
////
////            List<Address> addresses = addressService.getAddressesByUserId(userId);
////
////
////            return ResponseEntity.ok(addresses);
////        } catch (Exception e) {
////            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
////                    "code", 500,
////                    "message", "Lỗi khi lấy thông tin người dùng"
////            ));
////        }
////    }
////
////    @PostMapping("/add")
////    public ResponseEntity<?> getUserInfo(@RequestBody Address address) {
////        try {
////            addressService.save(address);
////            return ResponseEntity.ok(Map.of("data", address, "message", "Thêm địa chỉ mới thành công"));
////        } catch (Exception e) {
////            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
////                    "code", 500,
////                    "message", "Lỗi khi lấy thông tin người dùng"
////            ));
////        }
////    }
////
////    @PostMapping("/default")
////    public ResponseEntity<?> chooseAddressDefault(@RequestHeader("Authorization") String token, @RequestParam int addressId) {
////        try {
////            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
////            Users user = userService.getUserById(userId);
////            if (user == null) {
////                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
////                        "code", 404,
////                        "message", "Không tìm thấy người dùng"
////                ));
////            } else {
////                Address a = addressService.findAddressById(addressId).orElse(null);
////                if (a != null) {
////                    addressService.setDefaultAddress(userId, a.getId());
////                }
////                return ResponseEntity.ok(HttpStatus.OK);
////
////            }
////        } catch (Exception e) {
////            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
////                    "code", 500,
////                    "message", "Lỗi khi lấy thông tin người dùng"
////            ));
////        }
////    }
////
////
////}
