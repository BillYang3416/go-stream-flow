import { HttpErrorResponse } from '@angular/common/http';
import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTabGroup } from '@angular/material/tabs';
import { Router } from '@angular/router';
import { ApiPrefix, ApiService } from 'src/app/core/services/api.service';
import { environment } from 'src/environments/environment.dev';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit, AfterViewInit {
  @ViewChild(MatTabGroup) tabGroup!: MatTabGroup;

  signInForm!: FormGroup;

  signUpForm!: FormGroup;

  lineLoginUrl = environment.apiUrl + '/auth/line-login';

  constructor(
    private formBuilder: FormBuilder,
    private apiSvc: ApiService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.signInForm = this.formBuilder.group({
      username: ['', [Validators.required, Validators.minLength(4)]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });

    this.signUpForm = this.formBuilder.group({
      username: ['', [Validators.required, Validators.minLength(4)]],
      displayName: ['', [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  ngAfterViewInit(): void {
    this.tabGroup.selectedTabChange.subscribe((_) => {
      this.onReset();
    });
  }

  onSignIn(): void {
    if (this.signInForm.invalid) {
      this.snackBar.open('Please enter username and password correctly', 'OK');
      return;
    }
    const value = this.signInForm.value;
    this.apiSvc.post(ApiPrefix.AUTH, 'login', value).subscribe({
      next: (_) => {
        this.onReset();
        this.snackBar.open('Sign in successfully!', '', {
          duration: 2000,
        });
        this.router.navigate(['/file-flow']);
      },
      error: (err: HttpErrorResponse) => {
        this.snackBar.open(err.error.message, '', { duration: 2000 });
      },
    });
  }

  onSignUp(): void {
    if (this.signUpForm.invalid) {
      this.snackBar.open(
        'Please enter username,password and nickname correctly',
        'OK'
      );
      return;
    }
    const value = this.signUpForm.value;
    this.apiSvc.post(ApiPrefix.AUTH, 'register', value).subscribe({
      next: (_) => {
        this.onReset();
        this.snackBar.open('Sign up successfully, please try to sign in.', '', {
          duration: 2000,
        });
        this.tabGroup.selectedIndex = 0;
      },
      error: (err: HttpErrorResponse) => {
        this.snackBar.open(err.error.message, '', { duration: 2000 });
      },
    });
  }

  onReset() {
    this.signInForm.reset();
    this.signInForm.setErrors(null);
    this.signUpForm.reset();
    this.signUpForm.setErrors(null);
  }
}
