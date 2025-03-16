import 'package:car/api/api_service.dart';
import 'package:car/auth/login.dart';
import 'package:flutter/material.dart';

class SettingsPage extends StatefulWidget {
  @override
  _SettingsPageState createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final _apiService = ApiService();
  late Future<Map<String, dynamic>> _userFuture;
  final _formKey = GlobalKey<FormState>();
  bool _isEditingEmail = false;
  bool _isEditingPhone = false;
  late TextEditingController _emailController;
  late TextEditingController _phoneController;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _userFuture = _apiService.getUserInfo();
    _emailController = TextEditingController();
    _phoneController = TextEditingController();
  }

  @override
  void dispose() {
    _emailController.dispose();
    _phoneController.dispose();
    super.dispose();
  }

  Future<void> _saveChanges(String field) async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });
      try {
        await _apiService.updateUserInfo(
          _emailController.text,
          _phoneController.text,
        );
        setState(() {
          if (field == 'Email') {
            _isEditingEmail = false;
          } else if (field == 'Phone') {
            _isEditingPhone = false;
          }
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('$field updated successfully'),
            backgroundColor: Colors.green,
          ),
        );
        setState(() {
          _userFuture = _apiService.getUserInfo(); // Refresh user data
        });
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to update $field: $e'),
            backgroundColor: Colors.red,
          ),
        );
      } finally {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _logout() async {
    setState(() {
      _isLoading = true;
    });
    try {
      await _apiService.logout();
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => LoginPage()),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Logout failed: $e'),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Settings'),
      ),
      body: FutureBuilder<Map<String, dynamic>>(
        future: _userFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData) {
            return Center(child: Text('No user data found'));
          }

          final userData = snapshot.data!;
          _emailController.text = userData['email'];
          _phoneController.text = userData['phone_number'];

          return Padding(
            padding: EdgeInsets.all(16),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'My Information',
                    style: TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  SizedBox(height: 16),
                  Card(
                    elevation: 4,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Padding(
                      padding: EdgeInsets.all(16),
                      child: Column(
                        children: [
                          _buildInfoRow(
                              'Name', userData['name'], false, userData),
                          _buildInfoRow('National Number',
                              userData['national_number'], false, userData),
                          _buildInfoRow('Email', userData['email'], true,
                              userData, 'Email'),
                          _buildInfoRow('Phone', userData['phone_number'], true,
                              userData, 'Phone'),
                        ],
                      ),
                    ),
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
                          onPressed: _logout,
                          child: Text(
                            'Logout',
                            style: TextStyle(fontSize: 18),
                          ),
                        ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, bool isEditable,
      Map<String, dynamic> userData,
      [String? field]) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(fontSize: 16),
          ),
          if (isEditable &&
              ((label == 'Email' && _isEditingEmail) ||
                  (label == 'Phone' && _isEditingPhone)))
            Expanded(
              child: TextFormField(
                controller:
                    label == 'Email' ? _emailController : _phoneController,
                decoration: InputDecoration(
                  hintText: 'Enter $label',
                  border: OutlineInputBorder(),
                ),
                keyboardType: label == 'Phone'
                    ? TextInputType.phone
                    : TextInputType.emailAddress,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter $label';
                  }
                  if (label == 'Email' &&
                      !RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$')
                          .hasMatch(value)) {
                    return 'Please enter a valid email';
                  }
                  return null;
                },
              ),
            )
          else
            Text(
              value,
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
          if (isEditable &&
              !((label == 'Email' && _isEditingEmail) ||
                  (label == 'Phone' && _isEditingPhone)))
            IconButton(
              icon: Icon(Icons.edit),
              onPressed: () {
                setState(() {
                  if (label == 'Email') {
                    _isEditingEmail = true;
                  } else if (label == 'Phone') {
                    _isEditingPhone = true;
                  }
                });
              },
            ),
          if (isEditable &&
              ((label == 'Email' && _isEditingEmail) ||
                  (label == 'Phone' && _isEditingPhone)))
            Row(
              children: [
                IconButton(
                  icon: Icon(Icons.save),
                  onPressed: () => _saveChanges(field!),
                ),
                IconButton(
                  icon: Icon(Icons.cancel),
                  onPressed: () {
                    setState(() {
                      if (label == 'Email') {
                        _isEditingEmail = false;
                        _emailController.text = userData['email'];
                      } else if (label == 'Phone') {
                        _isEditingPhone = false;
                        _phoneController.text = userData['phone_number'];
                      }
                    });
                  },
                ),
              ],
            ),
        ],
      ),
    );
  }
}
