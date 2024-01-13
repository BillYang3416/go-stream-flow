import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ApiService } from 'src/app/core/services/api.service';
import { environment } from 'src/environments/environment.dev';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {
  signInForm!: FormGroup;

  signUpForm!: FormGroup;

  lineLoginUrl = environment.apiUrl + '/auth/line-login';

  constructor(
    private formBuilder: FormBuilder,
    private apiSvc: ApiService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.signInForm = this.formBuilder.group({
      username: ['', [Validators.required, Validators.minLength(4)]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });

    this.signUpForm = this.formBuilder.group({
      username: ['', [Validators.required, Validators.minLength(4)]],
      nickname: ['', [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  onSignIn(): void {
    if (this.signInForm.invalid) {
      this.snackBar.open('Please enter username and password correctly', 'OK');
      return;
    }
  }

  onSignUp(): void {
    if (this.signInForm.invalid) {
      this.snackBar.open(
        'Please enter username,password and nickname correctly',
        'OK'
      );
      return;
    }
  }
}
