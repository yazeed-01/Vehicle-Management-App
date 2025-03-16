import 'package:car/api/api_service.dart';
import 'package:car/auth/add_contact_info_page.dart';
import 'package:car/home_page.dart';
import 'package:flutter/material.dart';

class KeyVerificationPage extends StatefulWidget {
  final String nationalNumber;

  KeyVerificationPage({required this.nationalNumber});

  @override
  _KeyVerificationPageState createState() => _KeyVerificationPageState();
}

class _KeyVerificationPageState extends State<KeyVerificationPage> {
  final _formKey = GlobalKey<FormState>();
  String _specificKey = '';
  final _apiService = ApiService();
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Key Verification'),
      ),
      body: Container(
        padding: EdgeInsets.symmetric(horizontal: 24.0),
        child: Center(
          child: SingleChildScrollView(
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.vpn_key,
                    size: 80,
                    color: Theme.of(context).primaryColor,
                  ),
                  SizedBox(height: 16),
                  Text(
                    'Enter Verification Key',
                    style: TextStyle(
                      fontSize: 28,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  SizedBox(height: 48),
                  Card(
                    elevation: 8,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Padding(
                      padding: EdgeInsets.all(16),
                      child: Column(
                        children: [
                          Text(
                            'National Number: ${widget.nationalNumber}',
                            style: TextStyle(fontSize: 16),
                          ),
                          SizedBox(height: 16),
                          TextFormField(
                            decoration: InputDecoration(
                              labelText: 'Specific Key',
                              prefixIcon: Icon(Icons.key),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter the verification key';
                              }
                              return null;
                            },
                            onSaved: (value) {
                              _specificKey = value!;
                            },
                          ),
                          SizedBox(height: 24),
                          _isLoading
                              ? CircularProgressIndicator()
                              : ElevatedButton(
                                  style: ElevatedButton.styleFrom(
                                    minimumSize: Size(double.infinity, 50),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                  ),
                                  onPressed: () async {
                                    if (_formKey.currentState!.validate()) {
                                      _formKey.currentState!.save();
                                      setState(() {
                                        _isLoading = true;
                                      });
                                      try {
                                        final loginData =
                                            await _apiService.login(
                                                widget.nationalNumber,
                                                _specificKey);
                                        // Check if phone_number or email is missing
                                        final phoneNumber =
                                            loginData['phone_number']
                                                as String?;
                                        final email =
                                            loginData['email'] as String?;
                                        if (phoneNumber == null ||
                                            phoneNumber.isEmpty ||
                                            email == null ||
                                            email.isEmpty) {
                                          Navigator.pushReplacement(
                                            context,
                                            MaterialPageRoute(
                                                builder: (context) =>
                                                    AddContactInfoPage()),
                                          );
                                        } else {
                                          Navigator.pushReplacement(
                                            context,
                                            MaterialPageRoute(
                                                builder: (context) =>
                                                    HomePage()),
                                          );
                                        }
                                      } catch (e) {
                                        ScaffoldMessenger.of(context)
                                            .showSnackBar(
                                          SnackBar(
                                            content:
                                                Text('Verification failed: $e'),
                                            backgroundColor: Colors.red,
                                          ),
                                        );
                                      } finally {
                                        setState(() {
                                          _isLoading = false;
                                        });
                                      }
                                    }
                                  },
                                  child: Text(
                                    'Verify',
                                    style: TextStyle(fontSize: 18),
                                  ),
                                ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
