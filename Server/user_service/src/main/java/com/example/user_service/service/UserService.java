package com.example.user_service.service;

import com.example.user_service.dto.SignupCustomerRequest;
import com.example.user_service.dto.SignupEmployeeRequest;
import com.example.user_service.dto.UserResponseDTO;
import com.example.user_service.model.*;
import com.example.user_service.repository.*;
import jakarta.transaction.Transactional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.time.Duration;

@Service
public class UserService {

    @Autowired
    RoleRepository roleRepository;

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private CustomerRepository customerRepository;

    @Autowired
    private EmployeeRepository employeeRepository;

    @Autowired
    private DriverRepository driverRepository;
    @Autowired
    private RedisTemplate<String, String> redisTemplate;
    @Autowired
    private EmailService emailService;

    public void save(User u) {
        userRepository.save(u);
    }

    private PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();

    public UserService(PasswordEncoder passwordEncoder) {
        this.passwordEncoder = passwordEncoder;
    }

    public boolean authenticateUser(String username, String password) {

        Optional<User> user = userRepository.findUserByUsername(username);
        System.out.println(user.toString());
        return user.isPresent() && passwordEncoder.matches(password, user.get().getPassword());
    }

    public void updateUser(User user) {
        userRepository.save(user);
    }

    public User findByUsername(String username) {
        return userRepository.findUserByUsername(username)
                .orElseThrow(() -> new RuntimeException("Tài khoản không tồn tại: " + username));
    }

    public User findById(int id) {
        return userRepository.findById(id).get();
    }

    public boolean validateResetToken(String token) {
        if (userRepository.findByResetToken(token) != null) {
            return true;
        }
        return false;
    }

    public boolean updatePassword(String token, String newPassword) {
        User user = userRepository.findByResetToken(token);
        System.out.println(user.toString());
        if (user != null) {
            user.setPassword(passwordEncoder.encode(newPassword));
            user.setResetToken(null);
            userRepository.save(user);
            return true;
        }
        return false;
    }

    public void savePasswordResetToken(User user, String token) {
        user.setResetToken(token);
        userRepository.save(user);
    }

    public UserService() {
    }

    public UserService(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    private static final String OTP_PREFIX = "otp:";
    private static final Duration OTP_TTL = Duration.ofMinutes(5); // OTP hết hạn sau 5 phút

    public boolean existsByUsername(String username) {
        return userRepository.existsByUsername(username);
    }

    // Lưu OTP vào Redis với thời gian sống (TTL)
    public void saveOtp(String username, String otp) {
        String key = OTP_PREFIX + username;
        redisTemplate.opsForValue().set(key, otp, OTP_TTL);
    }

    // Xác thực OTP từ Redis
    public boolean verifyOtp(String username, String otp) {
        String key = OTP_PREFIX + username;
        String storedOtp = redisTemplate.opsForValue().get(key);
        return otp != null && otp.equals(storedOtp);
    }

    @Transactional
    public void registerUser(SignupCustomerRequest r) {
        System.out.println(r.getPassword());
        User newUser = new User();
        newUser.setUsername(r.getUsername());
        newUser.setPassword(passwordEncoder.encode(r.getPassword()));
        newUser.setCreatedAt(new Date());
        userRepository.save(newUser);

        Customer newCustomer = new Customer();
        newCustomer.setUser(newUser);
        newCustomer.setPhoneNumber(r.getPhoneNumber());
        newCustomer.setFullName(r.getFullName());
        newCustomer.setAddress(r.getAddress());
        newCustomer.setGender(r.getGender());
        customerRepository.save(newCustomer);

        Optional<Role> customerRoleOpt = Optional.ofNullable(roleRepository.findByRoleName("ROLE_CUSTOMER"));
        Role customerRole = customerRoleOpt.get();
        newUser.setRoles(Collections.singleton(customerRole));
        redisTemplate.delete(OTP_PREFIX + r.getUsername());
    }

    public String generateRandomPassword(int length) {
        String chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        Random rnd = new Random();
        StringBuilder sb = new StringBuilder(length);
        for (int i = 0; i < length; i++) {
            sb.append(chars.charAt(rnd.nextInt(chars.length())));
        }
        return sb.toString();
    }

    public void registerEmployee(SignupEmployeeRequest request) {
        User user = new User();
        user.setUsername(request.getUsername());
        String randomPassword = generateRandomPassword(8);
        user.setPassword(passwordEncoder.encode(randomPassword));
        user.setCreatedAt(new Date());
        String roleName = request.getEmployeeType().name();
        Role role = roleRepository.findByRoleName(roleName);
        if (role == null) {
            throw new RuntimeException("Role không tồn tại: " + roleName);
        }
        user.setRoles(Collections.singleton(role));

        userRepository.save(user);

        Employee employee = new Employee();
        employee.setUser(user);
        employee.setIdentityNumber(request.getIdentityNumber());
        employee.setHiredDate(new Date());
        employee.setFullName(request.getFullName());
        employee.setPhoneNumber(request.getPhoneNumber());
        employee.setAddress(request.getAddress());
        employee.setGender(request.getGender());
        employee.setDateOfBirth(request.getDateOfBirth());

        employeeRepository.save(employee);

        emailService.sendAccountInfoEmail(user.getUsername(), user.getUsername(), randomPassword);

        if (request.getEmployeeType() == EmployeeType.ROLE_DRIVER) {
            DriverInfo driverInfo = new DriverInfo();
            driverInfo.setEmployee(employee);
            driverInfo.setLicenseNumber(request.getLicenseNumber());
            driverInfo.setLicenseClass(request.getLicenseClass());
            driverInfo.setLicenseIssuedDate(request.getLicenseIssuedDate());
            driverInfo.setLicenseExpiryDate(request.getLicenseExpiryDate());
            driverInfo.setVehicleType(request.getVehicleType());
            driverRepository.save(driverInfo);
        }
    }

    public String getUserRole(String username) {

        Optional<User> user = userRepository.findUserByUsername(username);
        if (user.isPresent() && !user.get().getRoles().isEmpty()) {

            return user.get().getRoles().iterator().next().getRoleName();
        }
        return null;
    }

    public User getUserById(int id) {
        return userRepository.findById(id).orElse(null);
    }

    // public Integer getUserIdByEmail(String email) {
    // return userRepository.findIdByEmail(email).orElse(null);
    // }
    public List<UserResponseDTO> getUsersByRole(String roleName) {
        List<User> users = userRepository.findAll();
        List<UserResponseDTO> userDTOs = new ArrayList<>();

        for (User user : users) {
            boolean hasRole = user.getRoles().stream()
                    .anyMatch(role -> role.getRoleName().equals(roleName));

            if (hasRole) {
                UserResponseDTO dto = createUserDTO(user);
                if (dto != null) {
                    userDTOs.add(dto);
                }
            }
        }

        return userDTOs;
    }

    /**
     * Get all users with their details
     * 
     * @return list of all users
     */
    public List<UserResponseDTO> getAllUsers() {
        List<User> users = userRepository.findAll();
        List<UserResponseDTO> userDTOs = new ArrayList<>();

        for (User user : users) {
            UserResponseDTO dto = createUserDTO(user);
            if (dto != null) {
                userDTOs.add(dto);
            }
        }

        return userDTOs;
    }

    /**
     * Create a UserResponseDTO from a User entity
     * 
     * @param user the user entity
     * @return UserResponseDTO with appropriate details based on role
     */
    private UserResponseDTO createUserDTO(User user) {
        String role = getUserRole(user.getUsername());
        if (role == null) {
            return new UserResponseDTO(user);
        }

        if (role.equals("ROLE_CUSTOMER")) {
            Customer customer = customerRepository.findById(user.getId()).orElse(null);
            return UserResponseDTO.fromCustomer(user, customer);
        } else if (role.startsWith("ROLE_")) {
            // For employees (including drivers, receptionists, operators)
            Employee employee = employeeRepository.findById(user.getId()).orElse(null);
            return UserResponseDTO.fromEmployee(user, employee);
        }

        return new UserResponseDTO(user);
    }
}